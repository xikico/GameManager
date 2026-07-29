package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Game struct {
	Id          string    `gorm:"column:id;type:varchar(64);primaryKey;not null;" json:"id"`
	IconPath    string    `gorm:"column:icon_path;type:text;comment:图标;" json:"icon_path"`                         // 图标
	Name        string    `gorm:"column:name;type:varchar(50);comment:名称;not null;" json:"name"`                   // 名称
	NickName    string    `gorm:"column:nick_name;type:varchar(50);comment:昵称;" json:"nick_name"`                  // 昵称
	Series      string    `gorm:"column:series;type:varchar(100);comment:系列;" json:"series"`                       // 系列
	Description string    `gorm:"column:description;type:varchar(512);comment:描述;" json:"description"`             // 描述
	Path        string    `gorm:"column:path;type:varchar(1024);comment:游戏主路径;not null;" json:"path"`              // 游戏主路径
	StartPath   string    `gorm:"column:start_path;type:varchar(1024);comment:游戏启动路径;not null;" json:"start_path"` // 游戏启动路径
	CategoryId  string    `gorm:"column:category;type:varchar(64);comment:类别;" json:"category"`                    // 类别
	Imgs        []byte    `gorm:"column:imgs;type:blob;comment:图标;" json:"imgs"`                                   // 展示图
	IsPlay      bool      `gorm:"column:is_play;type:bool;default:false;comment:是否玩过;" json:"is_play"`             // 是否玩过
	IsDel       bool      `gorm:"column:is_del;type:bool;default:false;comment:是否删除;" json:"is_del"`               // 是否删除
	InsertTime  time.Time `gorm:"column:insert_time;type:datetime;comment:插入时间;" json:"insert_time"`               // 插入时间
}

func (Game) TableName() string {
	return "game" // 显式指定表名，覆盖默认的复数 games
}

// BeforeCreate GORM的钩子函数，在insert之前自动调用
func (g *Game) BeforeCreate(tx *gorm.DB) error {
	// 确保 UUID 为空时才生成（避免重复生成）
	if g.Id == "" {
		uuidObj, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		g.Id = strings.ReplaceAll(uuidObj.String(), "-", "")
	}
	if g.InsertTime.IsZero() {
		g.InsertTime = time.Now()
	}
	return nil
}
