package gin_support

import (
	"GameManager/adapters/in/gin_support/api"
	"GameManager/adapters/in/gin_support/middle"
	"GameManager/adapters/in/utils"
	"GameManager/domain/ports/in"
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
)

func initAPI(app *gin.Engine, manager in.GameManager) {
	game := app.Group("/game")
	ginAPI := api.NewGinAPI(manager)
	game.POST("/predict_game", ginAPI.PredictGame)
	game.POST("/delete_game", ginAPI.DeleteGame)
	game.GET("/get_game_by_category", ginAPI.GetGameByCategory)
	game.GET("/get_game_imgs", ginAPI.GetGameImgs)
	game.GET("/get_all_category", ginAPI.GetAllCategory)
	game.POST("/open_game", ginAPI.OpenGame)
	game.POST("/get_game_by_condition", ginAPI.GetGameByCondition)
	game.POST("/edit_game", ginAPI.EditGame)
	game.POST("/open_folder", api.OpenFolder)
}

func StartGinServer(manager in.GameManager) error {
	cfg := GetConfig()
	logger := utils.GetLogger()
	app := gin.Default()
	app.SetTrustedProxies([]string{cfg.Server.Host})
	app.Use(middle.GlobalErrorHandler(logger))
	initAPI(app, manager) // 注册路由一定要在注册中间件之后，不然中间件不会生效
	app.Static("/ui", "./ui")
	app.NoRoute(func(c *gin.Context) {
		indexFile := filepath.Join("ui", "index.html")
		c.File(indexFile) // Gin 的 File 方法会自动处理 MIME 和缓存
	})
	return app.Run(fmt.Sprintf(":%d", config.Server.Port))
}
