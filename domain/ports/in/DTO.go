package in

import (
	"time"
)

// 这里是直接复制的领域对象，因为输入端口不应该直接使用领域对象作为入参
// 理由是要明确边界、依赖倒置，领域对象的构建不应该依赖外部服务
// 说得实际一点：如果业务实体要变，且变的部分不需要用户输入，那么这样写就不需要更改外部服务（只改DTO转化为领域对象）

type GameDTO struct {
	Id          string
	IconPath    string      // 图标
	Name        string      // 名称
	NickName    string      // 昵称
	Series      string      // 系列
	Description string      // 描述
	Path        string      // 游戏主路径
	StartPath   string      // 游戏启动路径
	Category    CategoryDTO // 类别
	Imgs        []Img       // 展示图
	IsPlay      bool        // 是否玩过
	IsDel       bool        // 是否删除
	InsertTime  time.Time   // 插入时间
}

type Img []byte

type CategoryDTO struct {
	Id   string
	Name string // 类型名称
	Num  int    // 共有多少个游戏
}

type SearchGameConditionDTO struct {
	Name            string
	InsertTimeStart time.Time
	InsertTimeEnd   time.Time
	Description     string
	Series          string
	IsPlay          *bool
	CategoryDTO
}

type SettingsDTO struct {
	ClipboardImageDetectionEnabled bool
}

type ImportGameResult struct {
	IsMany   bool
	NeedPass bool
	Message  string
}
