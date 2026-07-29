package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const oldDBPath = `E:\PlayProgram\GameManegerEXE\game_manager.db`
const newDBPath = `./game_manager.db`

// ================= OLD =================

type OldGame struct {
	Id          string    `gorm:"column:id"`
	IconPath    string    `gorm:"column:icon_path"`
	Name        string    `gorm:"column:name"`
	NickName    string    `gorm:"column:nick_name"`
	Description string    `gorm:"column:description"`
	Path        string    `gorm:"column:path"`
	StartPath   string    `gorm:"column:start_path"`
	Category    string    `gorm:"column:category"`
	IsDel       bool      `gorm:"column:is_del"`
	InsertTime  time.Time `gorm:"column:insert_time"`
}

func (OldGame) TableName() string {
	return "game"
}

// ================= NEW =================

type Game struct {
	Id          string    `gorm:"column:id;primaryKey"`
	IconPath    string    `gorm:"column:icon_path"`
	Name        string    `gorm:"column:name"`
	NickName    string    `gorm:"column:nick_name"`
	Series      string    `gorm:"column:series"`
	Description string    `gorm:"column:description"`
	Path        string    `gorm:"column:path"`
	StartPath   string    `gorm:"column:start_path"`
	CategoryId  string    `gorm:"column:category"`
	Imgs        []byte    `gorm:"column:imgs"`
	IsPlay      bool      `gorm:"column:is_play"`
	IsDel       bool      `gorm:"column:is_del"`
	InsertTime  time.Time `gorm:"column:insert_time"`
}

func (Game) TableName() string {
	return "game"
}

func (g *Game) BeforeCreate(tx *gorm.DB) error {
	if g.Id == "" {
		u, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		g.Id = strings.ReplaceAll(u.String(), "-", "")
	}
	if g.InsertTime.IsZero() {
		g.InsertTime = time.Now()
	}
	return nil
}

type Category struct {
	Id   string `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:name"`
	Num  int    `gorm:"column:num"`
}

func (Category) TableName() string {
	return "category"
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.Id == "" {
		u, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		c.Id = strings.ReplaceAll(u.String(), "-", "")
	}
	return nil
}

func main() {
	// old db
	oldDB, err := gorm.Open(sqlite.Open(oldDBPath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// new db
	newDB, err := gorm.Open(sqlite.Open(newDBPath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// create tables
	if err := newDB.AutoMigrate(&Category{}, &Game{}); err != nil {
		log.Fatal(err)
	}

	var oldGames []OldGame
	if err := oldDB.Find(&oldGames).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Println("old games count:", len(oldGames))

	categoryMap := map[string]string{}
	categoryCount := map[string]int{}

	// collect categories
	for _, g := range oldGames {
		cat := strings.TrimSpace(g.Category)
		if cat == "" {
			cat = "未分类"
		}
		categoryCount[cat]++
	}

	// insert categories
	for catName, num := range categoryCount {
		c := Category{
			Name: catName,
			Num:  num,
		}

		if err := newDB.Create(&c).Error; err != nil {
			log.Fatal(err)
		}

		categoryMap[catName] = c.Id
	}

	// migrate games
	newGames := make([]Game, 0, len(oldGames))

	for _, oldGame := range oldGames {
		cat := strings.TrimSpace(oldGame.Category)
		if cat == "" {
			cat = "未分类"
		}

		newGame := Game{
			Id:          oldGame.Id,
			IconPath:    oldGame.IconPath,
			Name:        oldGame.Name,
			NickName:    oldGame.NickName,
			Series:      "", // old db doesn't have it
			Description: oldGame.Description,
			Path:        oldGame.Path,
			StartPath:   oldGame.StartPath,
			CategoryId:  categoryMap[cat],
			Imgs:        nil,   // default empty
			IsPlay:      false, // default
			IsDel:       oldGame.IsDel,
			InsertTime:  oldGame.InsertTime,
		}

		newGames = append(newGames, newGame)
	}

	if len(newGames) > 0 {
		if err := newDB.CreateInBatches(newGames, 100).Error; err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("migration success")
}
