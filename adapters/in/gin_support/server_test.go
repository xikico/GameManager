package gin_support

import (
	"GameManager/adapters/out/db/sqlite"
	unzip2 "GameManager/adapters/out/unzip"
	"GameManager/adapters/out/utils"
	"GameManager/domain/ports/in"
	"GameManager/domain/service"
	"testing"
)

func createGameManager() in.GameManager {
	db, _ := sqlite.GetDB()
	utils := utils.GetUtils() // 这是真的
	unzip := unzip2.Unzip     // 真的

	return service.NewGameManager(db, utils, unzip)
}

func TestApp(t *testing.T) {
	err := StartGinServer(createGameManager())
	if err != nil {
		t.Errorf("%e", err)
		return
	}
}
