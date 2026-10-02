package gin_support

import (
	"GameManager/adapters/out/db/fake"
	"GameManager/adapters/out/filesystem"
	unzip2 "GameManager/adapters/out/unzip"
	"GameManager/adapters/out/utils"
	"GameManager/domain/ports/in"
	"GameManager/domain/service"
	"testing"
)

func createGameManager() in.GameManager {
	db := fake.GetDB()
	utils := utils.GetUtils() // 这是真的
	unzip := unzip2.Unzip     // 真的

	return service.NewGameManager(db, utils, utils, utils, unzip, filesystem.NewInspector())
}

func TestApp(t *testing.T) {
	if createGameManager() == nil {
		t.Fatal("expected game manager")
	}
}
