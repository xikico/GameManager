package main

import (
	utils2 "GameManager/adapters/in/utils"
	"GameManager/adapters/in/wails_support"
	"GameManager/adapters/out/db/sqlite"
	unzip2 "GameManager/adapters/out/unzip"
	"GameManager/adapters/out/utils"
	"GameManager/domain/service"
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"go.uber.org/zap"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	log := utils2.GetLogger()
	db, err := sqlite.GetDB()
	if err != nil {
		log.Fatal("数据库初始化失败", zap.Error(err))
	}
	utils := utils.GetUtils()
	unzip := unzip2.Unzip
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gameManager := service.NewGameManager(db, utils, unzip)
	go utils.Img2Base64(ctx)

	app := wails_support.NewApp(gameManager)

	err = wails.Run(&options.App{
		Title:  "游戏管理器",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 59, A: 1},
		OnStartup:        app.Startup,
		OnShutdown:       app.Shutdown,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			// 使用相对路径 "./wails_data"
			WebviewUserDataPath: "./wails_data",
		},
	})
	if err != nil {
		log.Fatal("桌面应用启动失败", zap.Error(err))
	}
}
