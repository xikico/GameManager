package utils

import (
	"GameManager/domain/entity"
	"GameManager/domain/ports/in"
)

func GameToDTO(g entity.Game) in.GameDTO {
	// 转换图片切片
	imgs := make([]in.Img, len(g.Imgs))
	for i, data := range g.Imgs {
		imgs[i] = in.Img(data)
	}
	return in.GameDTO{
		Id:          g.Id,
		IconPath:    g.IconPath,
		Name:        g.Name,
		NickName:    g.NickName,
		Series:      g.Series,
		Description: g.Description,
		Path:        g.Path,
		StartPath:   g.StartPath,
		Category:    CategoryToDTO(g.Category),
		Imgs:        imgs,
		IsPlay:      g.IsPlay,
		IsDel:       g.IsDel,
		InsertTime:  g.InsertTime,
	}
}

func CategoryToDTO(c entity.Category) in.CategoryDTO {
	return in.CategoryDTO{
		Id:   c.Id,
		Name: c.Name,
		Num:  c.Num,
	}
}
