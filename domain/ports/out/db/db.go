package db

import "GameManager/domain/entity"

type DB interface {
	SaveGame(game entity.Game) error
	SaveGames(games []entity.Game) error
	GetGameByCategory(category entity.Category) ([]entity.Game, error)
	GetGameImgs(gameId string) ([]entity.Img, error) // 除了这个之外的GetGame都不返回图片（imgs）
	GetGameByCondition(condition entity.SearchGameCondition) ([]entity.Game, error)
	GetAllCategory() ([]entity.Category, error)
	EditGame(oldGame entity.Game, newGame entity.Game) error
}
