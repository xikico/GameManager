package api

import (
	"GameManager/adapters/in/gin_support/entity"
	"GameManager/adapters/in/utils"
	"GameManager/domain/ports/in"
	domainUtils "GameManager/domain/utils"
	"io"
	"os"

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
	path = domainUtils.NormalizePath(path)
	if !ok || path == "" {
		c.JSON(200, utils.Failure("", "缺少或无效的 path 参数"))
		return
	}
	isMany, confirmed := (*payload)["is_many"].(bool)
	if confirmed && isMany {
		if err := g.inPort.ConfirmImportMany(path); err != nil {
			c.JSON(200, utils.Failure("", err.Error()))
			return
		}
		c.JSON(200, utils.Success("", "已加入游戏"))
		return
	}
	password, _ := (*payload)["password"].(string)
	result, err := g.inPort.ImportGame(path, password)
	if err != nil {
		c.JSON(200, utils.Failure("", err.Error()))
		return
	}
	c.JSON(200, utils.Success(map[string]bool{"is_many": result.IsMany, "need_pass": result.NeedPass}, result.Message))
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
		return
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
	request.IconPath = domainUtils.NormalizePath(request.IconPath)
	request.Path = domainUtils.NormalizePath(request.Path)
	request.StartPath = domainUtils.NormalizePath(request.StartPath)
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
	// 未提交任何图片时保持原有截图不变（全量替换语义下需要回填现有图片）
	if len(imgsData) == 0 {
		existing, err := g.inPort.GetGameImgs(request.Id)
		if err == nil && existing != nil {
			imgsData = make([][]byte, len(existing))
			for i, img := range existing {
				imgsData[i] = img
			}
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
func (g *GinAPI) OpenFolder(c *gin.Context) {
	payload := getBody(c)
	if payload == nil {
		c.JSON(200, utils.Failure("", "缺少body参数"))
		return
	}
	path, ok := (*payload)["path"].(string)
	path = domainUtils.NormalizePath(path)
	if !ok {
		c.JSON(200, utils.Failure("", "缺少path参数"))
		return
	}
	if err := g.inPort.OpenFolder(path); err != nil {
		c.JSON(200, utils.Failure("", err.Error()))
		return
	}
	c.JSON(200, utils.Success("", "已打开文件夹"))
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
	path = domainUtils.NormalizePath(path)
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
