package service

import (
	"GameManager/adapters/out/db/fake"
	"GameManager/adapters/out/db/sqlite"
	unzip2 "GameManager/adapters/out/unzip"
	"GameManager/adapters/out/utils"
	"GameManager/domain/ports/in"
	"testing"
	"time"
)

func createFakeGameManager() *GameManager {
	db := fake.GetDB()
	utils := utils.GetUtils() // 这是真的
	unzip := unzip2.Unzip     // 真的

	return NewGameManager(db, utils, unzip)
}

func createGameManager() *GameManager {
	db, _ := sqlite.GetDB()
	utils := utils.GetUtils() // 这是真的
	unzip := unzip2.Unzip     // 真的

	return NewGameManager(db, utils, unzip)
}

func TestGameManager_PredictGame(t *testing.T) {
	gameManager := createFakeGameManager()
	game := in.GameDTO{
		Path: "G:\\Game\\slg\\BareBackReincarnation_Demo",
	}
	game, err := gameManager.PredictGame(game)
	if err != nil {
		t.Errorf("%e", err)
		return
	}
	t.Logf("%+v", game)
}

func TestGetGame(t *testing.T) {
	gameManager := createGameManager()
	games, err := gameManager.GetGameByCategory(in.CategoryDTO{
		Id:   "",
		Name: "galgame",
		Num:  0,
	})
	if err != nil {
		t.Errorf("按分类查询出错:%e", err)
		return
	}
	t.Logf("按分类查询：%+v", games)

	games, err = gameManager.GetGameByCondition(in.SearchGameConditionDTO{
		Name:          "魔女",
		InsertTimeEnd: time.Now(),
	})
	if err != nil {
		t.Errorf("按条件查询出错:%e", err)
		return
	}
	t.Logf("按条件查询：%+v", games)
}

func TestOpenGame(t *testing.T) {
	gameManager := createFakeGameManager()
	game := in.GameDTO{
		Path:       "E:\\utils",
		StartPath:  "E:\\utils\\witr.exe",
		Category:   in.CategoryDTO{},
		IsDel:      false,
		InsertTime: time.Now(),
	}
	err := gameManager.OpenGame(game)
	if err != nil {
		t.Errorf("打开游戏出错：%e", err)
		return
	}
}

func TestUnzip(t *testing.T) {
	gameManager := createFakeGameManager()
	s, err := gameManager.UnzipGame("E:\\utils\\testzip\\test.zip", "123")
	if err != nil {
		t.Errorf("解压出错：%e", err)
		return
	}
	t.Logf("解压成功:%s", s)
}
