package api

import (
	"GameManager/adapters/in/gin_support/entity"
	"GameManager/adapters/in/utils"
	"GameManager/domain/ports/in"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type GinAPI struct {
	inPort in.GameManager
}

func NewGinAPI(manager in.GameManager) GinAPI {
	return GinAPI{
		inPort: manager,
	}
}

// POST
func (g *GinAPI) PredictGame(c *gin.Context) {
	payload := getBody(c)
	if payload == nil {
		c.JSON(200, utils.Failure("", "缺少body参数"))
		return
	}
	path, ok := (*payload)["path"].(string)
	if !ok || path == "" {
		c.JSON(200, utils.Failure("", "缺少或无效的 path 参数"))
		return
	}
	isMany, ok := (*payload)["is_many"].(bool)
	password, passOk := (*payload)["password"].(string)
	if !ok {
		pathType, err := utils.AnalysisPath(path)
		if err != nil {
			c.JSON(200, utils.Failure("", err.Error()))
			return
		}
		switch pathType {
		case utils.SingleGameFolder:
			err = g.processSingleGame(path)
			if err != nil {
				c.JSON(200, utils.Failure("", err.Error()))
				return
			}
		case utils.ManyGameFolder:
			c.JSON(200, utils.Success(map[string]bool{"is_many": true, "need_pass": false}, "是否将其作为一个分类进行批量添加"))
			return
		case utils.ZipFile:
			var pass string
			if passOk {
				pass = password
			}
			game, err := g.inPort.UnzipGame(path, pass)
			if err != nil {
				if err.Error() == "密码错误或需要密码但未提供" {
					c.JSON(200, utils.Success(map[string]bool{"is_many": false, "need_pass": true}, "需要密码或密码错误"))
					return
				}
				c.JSON(200, utils.Failure("", err.Error()))
				return
			}
			err = g.processSingleGame(game)
			if err != nil {
				c.JSON(200, utils.Failure("", err.Error()))
				return
			}
		}

	} else if isMany {
		games := make([]in.GameDTO, 0, 1)
		f, err := os.Open(path)
		if err != nil {
			c.JSON(200, utils.Failure("", "打不开文件夹："+err.Error()))
			return
		}
		defer f.Close()

		entries, _ := f.ReadDir(-1) // -1 表示读取全部
		for _, entry := range entries {
			game := in.GameDTO{
				Path: filepath.Join(path, entry.Name()),
				Category: in.CategoryDTO{
					Name: filepath.Base(path),
				},
			}
			game, err = g.inPort.PredictGame(game)
			if err != nil {
				continue
			}
			games = append(games, game)
		}
		err = g.inPort.SaveGame(games)
		if err != nil {
			c.JSON(200, utils.Failure("", err.Error()))
			return
		}
	} else {
		err := g.processSingleGame(path)
		if err != nil {
			c.JSON(200, utils.Failure("", err.Error()))
			return
		}
	}
	c.JSON(200, utils.Success("", "已加入游戏"))
}

// GET
func (g *GinAPI) GetGameByCategory(c *gin.Context) {
	category_id := c.Query("category_id")
	if category_id == "" {
		c.JSON(200, utils.Failure("", "缺少 category_id 参数"))
		return
	}
	games, err := g.inPort.GetGameByCategory(in.CategoryDTO{Id: category_id})
	if err != nil {
		c.JSON(200, utils.Failure("", err.Error()))
		return
	}
	c.JSON(200, utils.Success(games, "查询成功"))
}

// GET
func (g *GinAPI) GetGameImgs(c *gin.Context) {
	gameId := c.Query("game_id")
	if gameId == "" {
		c.JSON(200, utils.Failure("", "缺少 game_id 参数"))
		return
	}
	imgs, err := g.inPort.GetGameImgs(gameId)
	if err != nil {
		c.JSON(200, utils.Failure("", err.Error()))
		return
	}
	imgsBase64 := make([]string, 0, len(imgs))
	for _, img := range imgs {
		base64, err := utils.Byte2Base64(img)
		if err != nil {
			continue
		}
		imgsBase64 = append(imgsBase64, base64)
	}
	c.JSON(200, utils.Success(imgsBase64, "查询成功"))
}

// POST
func (g *GinAPI) DeleteGame(c *gin.Context) {
	var request entity.GameRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	if request.Id == "" {
		c.JSON(200, utils.Failure("缺少参数", "缺少 category_id 参数"))
	}
	deleteAll := c.Query("delete_all") == "1"
	game := utils.GameVo2DTO(request, nil)
	err := g.inPort.DeleteGame(game, deleteAll)
	if err != nil {
		c.JSON(200, utils.Failure("", err.Error()))
		return
	}
	c.JSON(200, utils.Success("", "删除成功"))
}

// POST
func (g *GinAPI) OpenGame(c *gin.Context) {
	var request entity.GameRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	err := g.inPort.OpenGame(utils.GameVo2DTO(request, nil))
	if err != nil {
		c.JSON(200, utils.Failure("", err.Error()))
		return
	}
	c.JSON(200, utils.Success("", "已尝试打开游戏"))
}

// POST
func (g *GinAPI) GetGameByCondition(c *gin.Context) {
	var request entity.ConditionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	games, err := g.inPort.GetGameByCondition(utils.Condition2DTO(request))
	if err != nil {
		c.JSON(200, utils.Failure("", err.Error()))
		return
	}
	c.JSON(200, utils.Success(games, "查询成功"))
}

// GET
func (g *GinAPI) GetAllCategory(c *gin.Context) {
	allCategory, err := g.inPort.GetAllCategory()
	if err != nil {
		c.JSON(200, utils.Failure("", err.Error()))
		return
	}
	c.JSON(200, utils.Success(allCategory, "查询成功"))
}

// POST
func (g *GinAPI) EditGame(c *gin.Context) {
	var request entity.GameFormRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(200, utils.Failure("", "获取表单失败"))
		return
	}
	request.Path = strings.ReplaceAll(request.Path, "\"", "")
	request.StartPath = strings.ReplaceAll(request.StartPath, "\"", "")
	imgHeaders := form.File["imgs"]
	var imgsData [][]byte
	for _, header := range imgHeaders {
		src, err := header.Open()
		if err == nil {
			defer src.Close()
			data, _ := io.ReadAll(src)
			imgsData = append(imgsData, data)
		}
	}
	game := utils.GameFormVo2DTO(request, imgsData)
	if game.IconPath != "" {
		game.IconPath = procesIconPath(game.IconPath)
	}

	err = g.inPort.EditGame(game, game)
	if err != nil {
		c.JSON(200, utils.Failure("", err.Error()))
		return
	}
	c.JSON(200, utils.Success("", "修改成功"))
}

// POST
func OpenFolder(c *gin.Context) {
	payload := getBody(c)
	if payload == nil {
		c.JSON(200, utils.Failure("", "缺少body参数"))
		return
	}
	path, ok := (*payload)["path"].(string)
	if !ok {
		c.JSON(200, utils.Failure("", "缺少path参数"))
		return
	}
	cmd := exec.Command("explorer", path)
	cmd.Start()
}

func getBody(c *gin.Context) *map[string]interface{} {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(200, utils.Failure("", "请求格式错误"))
		c.Abort()
		return nil
	}
	return &payload
}

func procesIconPath(path string) string {
	if _, err := os.Stat(path); err != nil {
		return path
	} else {
		s, err := utils.Img2Base64(path)
		if err != nil {
			return ""
		}
		return s
	}
}

func (g *GinAPI) processSingleGame(path string) error {
	game := in.GameDTO{Path: path}
	var err error
	game, err = g.inPort.PredictGame(game)
	if err != nil {
		//c.JSON(200, utils.Failure("预测游戏信息出错：", err.Error()))
		return err
	}
	return g.inPort.SaveGame([]in.GameDTO{game})
}
