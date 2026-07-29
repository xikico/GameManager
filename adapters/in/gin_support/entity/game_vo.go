package entity

import "mime/multipart"

type GameRequest struct {
	Id           string `json:"id" binding:"omitempty,min=3,max=64"`
	IconPath     string `json:"icon_path" binding:"omitempty,min=3"`           // 图标
	Name         string `json:"name" binding:"required,min=1,max=500"`         // 名称
	NickName     string `json:"nick_name" binding:"omitempty,min=1,max=500"`   // 昵称
	Description  string `json:"description" binding:"omitempty,max=512"`       // 描述
	Path         string `json:"path" binding:"required,min=3,max=1024"`        // 游戏主路径
	StartPath    string `json:"start_path" binding:"omitempty,min=3,max=1024"` // 游戏启动路径
	CategoryName string `json:"category_name" binding:"omitempty,max=50"`      // 类别
	CategoryId   string `json:"category_id" binding:"omitempty,max=64"`        // 类别id
	Series       string `json:"series" binding:"omitempty,max=50"`             // 系列
	IsPlay       bool   `json:"is_play" binding:"omitempty""`                  // 是否玩过
}

type GameFormRequest struct {
	Id           string `form:"id"`
	IconPath     string `form:"icon_path"` // 可能是前端传的路径或标识
	Name         string `form:"name" binding:"required"`
	NickName     string `form:"nick_name"`
	Series       string `form:"series"`
	Description  string `form:"description"`
	Path         string `form:"path"`
	StartPath    string `form:"start_path"`
	CategoryId   string `form:"category_id"` // 用 ID 代替嵌套对象
	IsPlay       bool   `form:"is_play"`
	CategoryName string `form:"category_name" binding:"omitempty,max=50"` // 类别
}

type GameFiles struct {
	Imgs []*multipart.FileHeader // 多文件：展示图
}
