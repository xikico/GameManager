package sqlite

import (
	sqliteEntity "GameManager/adapters/out/db/sqlite/entity"
	"GameManager/adapters/out/db/sqlite/utils"
	"GameManager/domain/entity"
	"GameManager/domain/ports/out/db"
	"errors"
	"sync"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var once = sync.Once{}
var dbInstance *sqliteDB

func GetDB() (db.DB, error) {
	var err error
	if dbInstance == nil {
		once.Do(func() {
			dbInstance, err = newSqliteDB()
		})
	}
	return dbInstance, err
}

type sqliteDB struct {
	db          *gorm.DB
	allCategory []sqliteEntity.Category // 缓存，减少查询次数
}

func newSqliteDB() (*sqliteDB, error) {
	dsn := "file:./game_manager.db?_journal_mode=WAL&_busy_timeout=5000"

	// 连接 SQLite（纯 Go 驱动）
	dbObj, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		// 日志级别建议：开发用 Info，生产用 logger.Error 或 logger.Silent
		Logger: logger.Default.LogMode(logger.Error),
		// 推荐开启预编译语句缓存（对性能有帮助）
		PrepareStmt: true,

		// SQLite 默认事务行为通常不需要关闭，但可根据需求设置
		// SkipDefaultTransaction: true,
	})

	if err != nil {
		return nil, err
	}
	// 获取底层 *sql.DB，用于连接池调优（SQLite 下意义较小，但仍建议设置）
	sqlDB, err := dbObj.DB()
	if err != nil {
		return nil, err
	}

	// SQLite 连接池参数（防止极端情况下连接耗尽）
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(64) // SQLite 并发有限，建议不要太大
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(15 * time.Minute)

	// 自动建表
	if err := dbObj.AutoMigrate(&sqliteEntity.Category{}, &sqliteEntity.Game{}); err != nil {
		return nil, err
	}

	// 测试连接，顺手检查一下有没有默认分类
	category := sqliteEntity.Category{}
	if err = dbObj.Where("id = ?", "000000").Find(&category).Error; err != nil {
		return nil, err
	}
	if category.Name == "" {
		dbObj.Create(&sqliteEntity.Category{
			Id:   "000000",
			Name: "未分类",
			Num:  0,
		})
	}

	return &sqliteDB{
		db: dbObj,
	}, nil
}

func (s *sqliteDB) SaveGame(game entity.Game) error {
	db := s.db.Begin()
	categoryId := ""
	for _, category := range s.allCategory {
		if game.Category.Name == category.Name {
			categoryId = category.Id
			break
		}
	}
	dbGame := toDBGame(game)
	if dbGame.CategoryId != "" {
		if categoryId == "" {
			addCategory(db, toDBCategory(game.Category))
		} else {
			addCategoryNum(db, sqliteEntity.Category{Id: categoryId}, 1)
		}
	} else {
		dbGame.CategoryId = "000000"
		addCategoryNum(db, sqliteEntity.Category{Id: dbGame.CategoryId}, 1)
	}

	// 这里就是意思一下，其实以上任何一步出错都应该回滚
	err := db.Create(&dbGame).Error
	if err != nil {
		db.Rollback()
	} else {
		db.Commit()
	}
	return err
}

func (s *sqliteDB) SaveGames(games []entity.Game) error {
	db := s.db.Begin()
	currCategory := sqliteEntity.Category{}
	for _, category := range s.allCategory {
		if games[0].Category.Name == category.Name {
			currCategory = category
			break
		}
	}
	var newCategoryId string
	if currCategory.Id == "" {
		newCategoryId, _ = addCategory(db, sqliteEntity.Category{Name: games[0].Category.Name, Num: len(games)})
	} else {
		newCategoryId = currCategory.Id
		addCategoryNum(db, currCategory, len(games))
	}
	dbGames := make([]sqliteEntity.Game, 0, len(games))

	for _, g := range games {
		g.Category.Id = newCategoryId
		dbGames = append(dbGames, toDBGame(g))
	}

	err := db.Create(&dbGames).Error
	if err != nil {
		db.Rollback()
	} else {
		db.Commit()
	}
	return err
}

func (s *sqliteDB) GetGameByCategory(category entity.Category) ([]entity.Game, error) {
	if category.Id == "" && category.Name == "" {
		return nil, errors.New("参数不能为空")
	}
	var categoryDB sqliteEntity.Category
	if category.Id == "" {
		s.db.Select("id").Where("name = ?", category.Name).Find(&categoryDB)
		category.Id = categoryDB.Id
	}

	var games []sqliteEntity.Game
	err := s.db.
		Select("id", "icon_path", "name", "nick_name", "series",
			"description", "path", "start_path", "category",
			"is_play", "is_del", "insert_time").
		Where("category = ? AND is_del = ?", category.Id, false).
		Find(&games).Error

	if err != nil {
		return nil, err
	}

	result := make([]entity.Game, 0, len(games))
	for _, g := range games {
		result = append(result, toBizGame(g, category))
	}

	return result, nil
}

func (s *sqliteDB) GetGameImgs(gameId string) ([]entity.Img, error) {
	var game sqliteEntity.Game

	err := s.db.
		Select("imgs").
		Where("id = ?", gameId).
		First(&game).Error

	if err != nil {
		return nil, err
	}

	if game.Imgs == nil {
		return nil, nil
	}

	return utils.Unpack([]byte(game.Imgs))
}

func (s *sqliteDB) GetGameByCondition(cond entity.SearchGameCondition) ([]entity.Game, error) {
	query := s.db.
		Table("game g").
		Select(`
			g.id, g.icon_path, g.name, g.nick_name, g.series,
			g.description, g.path, g.start_path, g.category,
			g.is_play, g.is_del, g.insert_time,
			c.name as category_name
		`).
		Joins("LEFT JOIN category c ON g.category = c.id")

	query = query.Where("g.is_del = ?", false)

	if cond.Name != "" {
		query = query.Where("g.name LIKE ?", "%"+cond.Name+"%")
	}

	if cond.Description != "" {
		query = query.Where("g.description LIKE ?", "%"+cond.Description+"%")
	}

	if cond.Series != "" {
		query = query.Where("g.series LIKE ?", "%"+cond.Series+"%")
	}

	if !cond.InsertTimeStart.IsZero() {
		query = query.Where("g.insert_time >= ?", cond.InsertTimeStart)
	}

	if !cond.InsertTimeEnd.IsZero() {
		query = query.Where("g.insert_time <= ?", cond.InsertTimeEnd)
	}

	if cond.Category.Id != "" {
		query = query.Where("g.category = ?", cond.Category.Id)
	}

	if cond.IsPlay != nil {
		query = query.Where("g.is_play = ?", *cond.IsPlay)
	}

	var games []gameWithCategory
	if err := query.Scan(&games).Error; err != nil {
		return nil, err
	}

	result := make([]entity.Game, 0, len(games))
	for _, g := range games {
		result = append(result, entity.Game{
			Id:          g.Id,
			IconPath:    g.IconPath,
			Name:        g.Name,
			NickName:    g.NickName,
			Series:      g.Series,
			Description: g.Description,
			Path:        g.Path,
			StartPath:   g.StartPath,
			Category: entity.Category{
				Id:   g.CategoryId,
				Name: g.CategoryName,
			},
			IsPlay:     g.IsPlay,
			IsDel:      g.IsDel,
			InsertTime: g.InsertTime,
		})
	}

	return result, nil
}

func (s *sqliteDB) GetAllCategory() ([]entity.Category, error) {
	var categories []sqliteEntity.Category

	err := s.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	result := make([]entity.Category, 0, len(categories))
	for _, c := range categories {
		result = append(result, entity.Category{
			Id:   c.Id,
			Name: c.Name,
			Num:  c.Num,
		})
	}
	s.allCategory = categories

	return result, nil
}

func (s *sqliteDB) EditGame(oldGame entity.Game, newGame entity.Game) error {
	db := s.db.Begin()
	oldCategoryId := oldGame.Category.Id
	db.Model(&sqliteEntity.Game{}).Where("id = ?", oldGame.Id).Pluck("category", &oldCategoryId)
	update := map[string]interface{}{
		"icon_path":   newGame.IconPath,
		"name":        newGame.Name,
		"nick_name":   newGame.NickName,
		"series":      newGame.Series,
		"description": newGame.Description,
		"path":        newGame.Path,
		"start_path":  newGame.StartPath,
		"category":    newGame.Category.Id,
		"is_play":     newGame.IsPlay,
		"is_del":      newGame.IsDel,
	}

	// imgs 为完整最终列表：nil 表示不修改，空列表表示清空全部截图
	if newGame.Imgs != nil {
		update["imgs"] = string(utils.Pack(newGame.Imgs...))
	}

	if newGame.IsDel {
		subCategoryNum(db, sqliteEntity.Category{Id: newGame.Category.Id}, 1)
	} else if oldCategoryId != "" && oldCategoryId != newGame.Category.Id {
		subCategoryNum(db, sqliteEntity.Category{Id: oldCategoryId}, 1)
		addCategoryNum(db, sqliteEntity.Category{Id: newGame.Category.Id}, 1)
	}

	err := db.Model(&sqliteEntity.Game{}).
		Where("id = ?", oldGame.Id).
		Updates(update).Error
	if err != nil {
		db.Rollback()
	} else {
		db.Commit()
	}
	return err
}

func addCategory(db *gorm.DB, category sqliteEntity.Category) (string, error) {
	err := db.Create(&category).Error
	return category.Id, err
}

func addCategoryNum(db *gorm.DB, category sqliteEntity.Category, num int) error {
	return db.Model(&category).Where("id = ?", category.Id).Update("num", gorm.Expr("num + ?", num)).Error
}
func subCategoryNum(db *gorm.DB, category sqliteEntity.Category, num int) error {
	return db.Model(&category).Where("id = ?", category.Id).Update("num", gorm.Expr("num - ?", num)).Error
}

type gameWithCategory struct {
	sqliteEntity.Game
	CategoryName string `gorm:"column:category_name"`
}

func toDBGame(g entity.Game) sqliteEntity.Game {
	var imgs []byte
	if len(g.Imgs) > 0 {
		imgs = utils.Pack(g.Imgs...)
	}

	return sqliteEntity.Game{
		Id:          g.Id,
		IconPath:    g.IconPath,
		Name:        g.Name,
		NickName:    g.NickName,
		Series:      g.Series,
		Description: g.Description,
		Path:        g.Path,
		StartPath:   g.StartPath,
		CategoryId:  g.Category.Id,
		Imgs:        imgs,
		IsPlay:      g.IsPlay,
		IsDel:       g.IsDel,
		InsertTime:  g.InsertTime,
	}
}

func toDBCategory(c entity.Category) sqliteEntity.Category {
	return sqliteEntity.Category{
		Name: c.Name,
	}
}

func toBizGame(g sqliteEntity.Game, category entity.Category) entity.Game {
	return entity.Game{
		Id:          g.Id,
		IconPath:    g.IconPath,
		Name:        g.Name,
		NickName:    g.NickName,
		Series:      g.Series,
		Description: g.Description,
		Path:        g.Path,
		StartPath:   g.StartPath,
		Category:    category,
		IsPlay:      g.IsPlay,
		IsDel:       g.IsDel,
		InsertTime:  g.InsertTime,
	}
}
