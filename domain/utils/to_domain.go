package utils

import (
	"GameManager/domain/entity"
	"GameManager/domain/ports/in"
)

func GameToDomain(g in.GameDTO) *entity.Game {
	// 转换图片切片（保留 nil 语义：nil 表示未提供，空切片表示清空）
	var imgs []entity.Img
	if g.Imgs != nil {
		imgs = make([]entity.Img, len(g.Imgs))
		for i, data := range g.Imgs {
			imgs[i] = entity.Img(data)
		}
	}

	return &entity.Game{
		Id:          g.Id,
		IconPath:    g.IconPath,
		Name:        g.Name,
		NickName:    g.NickName,
		Series:      g.Series,
		Description: g.Description,
		Path:        g.Path,
		StartPath:   g.StartPath,
		Category:    *CategoryToDomain(g.Category),
		Imgs:        imgs,
		IsPlay:      g.IsPlay,
		IsDel:       g.IsDel,
		InsertTime:  g.InsertTime,
	}
}

func CategoryToDomain(c in.CategoryDTO) *entity.Category {
	return &entity.Category{
		Id:   c.Id,
		Name: c.Name,
		Num:  c.Num,
	}
}

func SearchGameConditionToDomain(s in.SearchGameConditionDTO) *entity.SearchGameCondition {
	return &entity.SearchGameCondition{
		Name:            s.Name,
		Series:          s.Series,
		Description:     s.Description,
		Category:        *CategoryToDomain(s.CategoryDTO),
		IsPlay:          s.IsPlay,
		InsertTimeStart: s.InsertTimeStart,
		InsertTimeEnd:   s.InsertTimeEnd,
	}
}
