package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type Category struct {
	Id   string `gorm:"column:id;type:varchar(64);primaryKey;not null;" json:"id"`
	Name string `gorm:"column:name;type:varchar(50);comment:名称;not null;" json:"name"`
	Num  int    `gorm:"column:num;type:Integer;comment:数量;" json:"num"`
}

func (Category) TableName() string {
	return "category" // 显式指定表名，覆盖默认的复数 games
}

// BeforeCreate GORM的钩子函数，在insert之前自动调用
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	// 确保 UUID 为空时才生成（避免重复生成）
	if c.Id == "" {
		uuidObj, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		c.Id = strings.ReplaceAll(uuidObj.String(), "-", "")
	}
	if c.Num == 0 && c.Name != "未分类" {
		c.Num = 1
	}
	return nil
}
