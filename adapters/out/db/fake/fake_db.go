package fake

import (
	"GameManager/domain/entity"
	"GameManager/domain/ports/out/db"
	"fmt"
)

func GetDB() db.DB {
	return &fakeDBImpl{settings: entity.Settings{ClipboardImageDetectionEnabled: true}}
}

type fakeDBImpl struct {
	settings entity.Settings
}

func (f *fakeDBImpl) SaveGame(game entity.Game) error {
	fmt.Printf("保存game：%+v \n", game)
	return nil
}

func (f *fakeDBImpl) SaveGames(games []entity.Game) error {
	fmt.Printf("保存games：%+v \n", games)
	return nil
}

func (f *fakeDBImpl) GetGameByCategory(category entity.Category) ([]entity.Game, error) {
	fmt.Printf("查询gameByCategory：%+v \n", category)
	return nil, nil
}

func (f *fakeDBImpl) GetGameImgs(gameId string) ([]entity.Img, error) {
	return nil, nil
}

func (f *fakeDBImpl) GetGameByCondition(condition entity.SearchGameCondition) ([]entity.Game, error) {
	return nil, nil
}

func (f *fakeDBImpl) GetAllCategory() ([]entity.Category, error) {
	return nil, nil
}

func (f *fakeDBImpl) EditGame(oldGame entity.Game, newGame entity.Game) error {
	return nil
}

func (f *fakeDBImpl) GetSettings() (entity.Settings, error) {
	return f.settings, nil
}

func (f *fakeDBImpl) SaveSettings(settings entity.Settings) error {
	f.settings = settings
	return nil
}
