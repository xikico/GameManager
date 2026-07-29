package main

import (
	"GameManager/adapters/in/gin_support"
	utils2 "GameManager/adapters/in/utils"
	"GameManager/adapters/out/db/sqlite"
	unzip2 "GameManager/adapters/out/unzip"
	"GameManager/adapters/out/utils"
	"GameManager/domain/service"
	"context"
	"go.uber.org/zap"
	"os/exec"
	"time"
)

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

	go func() {
		time.Sleep(300 * time.Millisecond) // 等待服务就绪
		cmd := exec.Command("cmd", "/c", "start", "http://127.0.0.1:10086/ui")
		err := cmd.Start()
		if err != nil {
			log.Fatal("自动打开浏览器失败:", zap.Error(err))
		}
	}()
	go utils.Img2Base64(ctx)

	err = gin_support.StartGinServer(gameManager)
	if err != nil {
		log.Fatal("服务器启动失败", zap.Error(err))
	}
}
