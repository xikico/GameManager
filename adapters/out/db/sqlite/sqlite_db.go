package sqlite

import (
	sqliteEntity "GameManager/adapters/out/db/sqlite/entity"
	"GameManager/adapters/out/db/sqlite/utils"
	"GameManager/domain/entity"
	"GameManager/domain/ports/out/db"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	once       sync.Once
	dbInstance *sqliteDB
	dbInitErr  error
)

func GetDB() (db.DB, error) {
	once.Do(func() {
		dbInstance, dbInitErr = newSqliteDB()
	})
	return dbInstance, dbInitErr
}

type sqliteDB struct {
	db *gorm.DB
}

func (s *sqliteDB) close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func newSqliteDB() (*sqliteDB, error) {
	return openSqliteDB("file:./game_manager.db?_journal_mode=WAL&_busy_timeout=5000")
}

func openSqliteDB(dsn string) (*sqliteDB, error) {
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
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(15 * time.Minute)

	// 自动建表
	if err := dbObj.AutoMigrate(&sqliteEntity.Category{}, &sqliteEntity.Game{}, &sqliteEntity.Settings{}); err != nil {
		return nil, err
	}

	// 测试连接，顺手检查一下有没有默认分类
	category := sqliteEntity.Category{}
	if err = dbObj.Where("id = ?", "000000").Find(&category).Error; err != nil {
		return nil, err
	}
	if category.Name == "" {
		if err := dbObj.Create(&sqliteEntity.Category{
			Id:   "000000",
			Name: "未分类",
			Num:  0,
		}).Error; err != nil {
			return nil, err
		}
	}
	var settings sqliteEntity.Settings
	if err = dbObj.First(&settings, 1).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		if err = dbObj.Create(&sqliteEntity.Settings{
			Id:                             1,
			ClipboardImageDetectionEnabled: true,
		}).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &sqliteDB{
		db: dbObj,
	}, nil
}

func (s *sqliteDB) GetSettings() (entity.Settings, error) {
	var settings sqliteEntity.Settings
	if err := s.db.First(&settings, 1).Error; err != nil {
		return entity.Settings{}, err
	}
	return entity.Settings{ClipboardImageDetectionEnabled: settings.ClipboardImageDetectionEnabled}, nil
}

func (s *sqliteDB) SaveSettings(settings entity.Settings) error {
	result := s.db.Model(&sqliteEntity.Settings{}).
		Where("id = ?", 1).
		Update("clipboard_image_detection_enabled", settings.ClipboardImageDetectionEnabled)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("设置记录不存在")
	}
	return nil
}

func (s *sqliteDB) SaveGame(game entity.Game) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		category, err := resolveCategory(tx, game.Category)
		if err != nil {
			return err
		}
		dbGame := toDBGame(game)
		dbGame.CategoryId = category.Id
		if err := tx.Create(&dbGame).Error; err != nil {
			return err
		}
		return changeCategoryCount(tx, category.Id, 1)
	})
}

func (s *sqliteDB) SaveGames(games []entity.Game) error {
	if len(games) == 0 {
		return errors.New("批量保存的游戏不能为空")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		category, err := resolveCategory(tx, games[0].Category)
		if err != nil {
			return err
		}
		dbGames := make([]sqliteEntity.Game, 0, len(games))
		for _, game := range games {
			if game.Category.Name != "" && !strings.EqualFold(game.Category.Name, games[0].Category.Name) {
				return errors.New("批量保存的游戏必须属于同一分类")
			}
			dbGame := toDBGame(game)
			dbGame.CategoryId = category.Id
			dbGames = append(dbGames, dbGame)
		}
		if err := tx.Create(&dbGames).Error; err != nil {
			return err
		}
		return changeCategoryCount(tx, category.Id, len(dbGames))
	})
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
	return result, nil
}

func (s *sqliteDB) EditGame(oldGame entity.Game, newGame entity.Game) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var persisted sqliteEntity.Game
		if err := tx.Select("id", "category", "is_del").Where("id = ?", oldGame.Id).First(&persisted).Error; err != nil {
			return err
		}

		categoryId := persisted.CategoryId
		if !newGame.IsDel || newGame.Category.Id != "" || newGame.Category.Name != "" {
			category, err := resolveCategory(tx, newGame.Category)
			if err != nil {
				return err
			}
			categoryId = category.Id
		}

		if err := applyCategoryTransition(tx, persisted.CategoryId, categoryId, persisted.IsDel, newGame.IsDel); err != nil {
			return err
		}

		update := map[string]interface{}{
			"icon_path": newGame.IconPath, "name": newGame.Name, "nick_name": newGame.NickName,
			"series": newGame.Series, "description": newGame.Description, "path": newGame.Path,
			"start_path": newGame.StartPath, "category": categoryId, "is_play": newGame.IsPlay, "is_del": newGame.IsDel,
		}
		if newGame.Imgs != nil {
			update["imgs"] = string(utils.Pack(newGame.Imgs...))
		}
		result := tx.Model(&sqliteEntity.Game{}).Where("id = ?", oldGame.Id).Updates(update)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func newCategoryId() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}

func resolveCategory(tx *gorm.DB, category entity.Category) (sqliteEntity.Category, error) {
	var result sqliteEntity.Category
	if category.Id != "" {
		if err := tx.Where("id = ?", category.Id).First(&result).Error; err == nil {
			return result, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return result, err
		}
	}
	name := strings.TrimSpace(category.Name)
	if name == "" {
		if err := tx.Where("id = ?", "000000").First(&result).Error; err != nil {
			return result, err
		}
		return result, nil
	}
	if err := tx.Where("name = ?", name).First(&result).Error; err == nil {
		return result, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return result, err
	}
	result = sqliteEntity.Category{Id: newCategoryId(), Name: name, Num: 0}
	if err := tx.Table(result.TableName()).Create(map[string]interface{}{"id": result.Id, "name": result.Name, "num": 0}).Error; err != nil {
		return result, err
	}
	return result, nil
}

func applyCategoryTransition(tx *gorm.DB, oldId, newId string, wasDeleted, isDeleted bool) error {
	switch {
	case !wasDeleted && isDeleted:
		return changeCategoryCount(tx, oldId, -1)
	case wasDeleted && !isDeleted:
		return changeCategoryCount(tx, newId, 1)
	case !wasDeleted && !isDeleted && oldId != newId:
		if err := changeCategoryCount(tx, oldId, -1); err != nil {
			return err
		}
		return changeCategoryCount(tx, newId, 1)
	default:
		return nil
	}
}

func changeCategoryCount(tx *gorm.DB, id string, delta int) error {
	result := tx.Model(&sqliteEntity.Category{}).Where("id = ?", id).
		Update("num", gorm.Expr("MAX(num + ?, 0)", delta))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("分类不存在: %s", id)
	}
	return nil
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
