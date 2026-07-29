package utils

import (
	"GameManager/adapters/in/gin_support/entity"
	"GameManager/domain/ports/in"
	"time"
)

func GameVo2DTO(vo entity.GameRequest, imgs [][]byte) in.GameDTO {
	// 转换图片数据 []Img
	imgList := make([]in.Img, len(imgs))
	for i, b := range imgs {
		imgList[i] = b
	}

	// 构造 CategoryDTO
	categoryDTO := in.CategoryDTO{
		Id:   vo.CategoryId,
		Name: vo.CategoryName,
	}

	return in.GameDTO{
		Id:          vo.Id,
		IconPath:    vo.IconPath,
		Name:        vo.Name,
		NickName:    vo.NickName,
		Series:      vo.Series,
		Description: vo.Description,
		Path:        vo.Path,
		StartPath:   vo.StartPath,
		Category:    categoryDTO,
		Imgs:        imgList,
		IsPlay:      vo.IsPlay,
		IsDel:       false,      // 新记录默认未删除
		InsertTime:  time.Now(), // 当前时间作为插入时间
	}
}

func GameFormVo2DTO(vo entity.GameFormRequest, imgs [][]byte) in.GameDTO {
	// 转换图片数据 []Img
	imgList := make([]in.Img, len(imgs))
	for i, b := range imgs {
		imgList[i] = b
	}

	// 构造 CategoryDTO
	categoryDTO := in.CategoryDTO{
		Id:   vo.CategoryId,
		Name: vo.CategoryName,
	}

	return in.GameDTO{
		Id:          vo.Id,
		IconPath:    vo.IconPath,
		Name:        vo.Name,
		NickName:    vo.NickName,
		Series:      vo.Series,
		Description: vo.Description,
		Path:        vo.Path,
		StartPath:   vo.StartPath,
		Category:    categoryDTO,
		Imgs:        imgList,
		IsPlay:      vo.IsPlay,
		IsDel:       false,      // 新记录默认未删除
		InsertTime:  time.Now(), // 当前时间作为插入时间
	}
}

func Condition2DTO(vo entity.ConditionRequest) in.SearchGameConditionDTO {
	return in.SearchGameConditionDTO{
		Name:            vo.Name,
		InsertTimeStart: vo.InsertTimeStart,
		InsertTimeEnd:   vo.InsertTimeEnd,
		Description:     vo.Description,
		Series:          vo.Series,
		IsPlay:          vo.IsPlay,
		CategoryDTO:     in.CategoryDTO{Id: vo.CategoryId},
	}
}
