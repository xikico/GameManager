package entity

import (
	"time"
)

type Img []byte

type Game struct {
	Id          string
	IconPath    string    // 图标
	Name        string    // 名称
	NickName    string    // 昵称
	Series      string    // 系列
	Description string    // 描述
	Path        string    // 游戏主路径
	StartPath   string    // 游戏启动路径
	Category    Category  // 类别
	Imgs        []Img     // 展示图
	IsPlay      bool      // 是否玩过
	IsDel       bool      // 是否删除
	InsertTime  time.Time // 插入时间
}

type SearchGameCondition struct {
	Name            string
	InsertTimeStart time.Time
	InsertTimeEnd   time.Time
	Description     string
	Series          string
	IsPlay          *bool
	Category
}
