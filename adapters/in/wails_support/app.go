// Package wails_support 提供桌面端入口适配器。
//
// 该适配器是六边形架构的 in 侧新端口实现：
// 只负责把 Wails 前端的方法调用翻译成 in.GameManager 接口调用，
// 自身不包含任何业务逻辑，核心业务全部复用 domain/service 层。
package wails_support

import (
	"GameManager/domain/ports/in"
	domainUtils "GameManager/domain/utils"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// App 是所有暴露给 Wails 前端的方法的载体。
type App struct {
	manager in.GameManager
	ctx     context.Context
}

func NewApp(manager in.GameManager) *App {
	return &App{manager: manager}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Shutdown(ctx context.Context) {}

// ---------- 分类 ----------

func (a *App) GetAllCategory() ([]in.CategoryDTO, error) {
	return a.manager.GetAllCategory()
}

// ---------- 设置 ----------

func (a *App) GetSettings() (in.SettingsDTO, error) {
	return a.manager.GetSettings()
}

func (a *App) UpdateSettings(settings in.SettingsDTO) error {
	return a.manager.UpdateSettings(settings)
}

// ---------- 游戏查询 ----------

// GetGames 按条件查询游戏（不传条件则返回全部）
func (a *App) GetGames(condition in.SearchGameConditionDTO) ([]in.GameDTO, error) {
	return a.manager.GetGameByCondition(condition)
}

func (a *App) GetGameByCategory(category in.CategoryDTO) ([]in.GameDTO, error) {
	return a.manager.GetGameByCategory(category)
}

// GetGameImgs 返回游戏展示图（base64 data URL 形式，可直接用于 <img src>）
func (a *App) GetGameImgs(gameId string) ([]string, error) {
	imgs, err := a.manager.GetGameImgs(gameId)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(imgs))
	for _, img := range imgs {
		result = append(result, dataURL(img))
	}
	return result, nil
}

// IconToBase64 将图标路径转换为 base64 data URL。
// 若路径不存在或读取失败返回空字符串。
func (a *App) IconToBase64(iconPath string) string {
	iconPath = domainUtils.NormalizePath(iconPath)
	if strings.HasPrefix(iconPath, "data:") {
		return iconPath
	}
	if iconPath == "" {
		return ""
	}
	data, err := os.ReadFile(iconPath)
	if err != nil {
		return ""
	}
	return dataURL(data)
}

// ---------- 游戏操作 ----------

func (a *App) OpenGame(game in.GameDTO) error {
	return a.manager.OpenGame(game)
}

func (a *App) DeleteGame(game in.GameDTO, delAll bool) error {
	return a.manager.DeleteGame(game, delAll)
}

// OpenFolder 在资源管理器中打开指定目录
func (a *App) OpenFolder(path string) error {
	return a.manager.OpenFolder(path)
}

// ---------- 添加游戏 ----------

type AddGameResult struct {
	IsMany   bool   `json:"is_many"`
	NeedPass bool   `json:"need_pass"`
	Message  string `json:"message"`
}

// AddGame 添加单个游戏：
// 文件夹 → 预测信息后入库；压缩包 → 先解压再入库；
// 包含多个子游戏的文件夹 → 返回 IsMany=true，由前端确认后调用 ConfirmAddMany。
func (a *App) AddGame(path, password string) (AddGameResult, error) {
	result, err := a.manager.ImportGame(path, password)
	return AddGameResult{IsMany: result.IsMany, NeedPass: result.NeedPass, Message: result.Message}, err
}

// ConfirmAddMany 将文件夹内的所有子目录批量添加，并把文件夹名作为分类
func (a *App) ConfirmAddMany(path string) error {
	return a.manager.ConfirmImportMany(path)
}

// PredictGame 预测游戏启动路径等信息（不保存）
func (a *App) PredictGame(game in.GameDTO) (in.GameDTO, error) {
	return a.manager.PredictGame(game)
}

// ---------- 编辑游戏 ----------

// EditGamePayload 编辑游戏的请求参数（与原 gin 表单接口等价，改用 JSON 传输）
type EditGamePayload struct {
	Id           string
	IconPath     string
	Name         string
	NickName     string
	Series       string
	Description  string
	Path         string
	StartPath    string
	CategoryId   string
	CategoryName string
	IsPlay       bool
	// Imgs 完整的最终展示图列表（data URL），全量替换；传空数组则清空截图
	Imgs []string
}

func (a *App) EditGame(payload EditGamePayload) error {
	game := in.GameDTO{
		Id:          payload.Id,
		IconPath:    domainUtils.NormalizePath(payload.IconPath),
		Name:        payload.Name,
		NickName:    payload.NickName,
		Series:      payload.Series,
		Description: payload.Description,
		Path:        domainUtils.NormalizePath(payload.Path),
		StartPath:   domainUtils.NormalizePath(payload.StartPath),
		Category:    in.CategoryDTO{Id: payload.CategoryId, Name: payload.CategoryName},
		IsPlay:      payload.IsPlay,
	}
	// Imgs != nil 时执行全量替换（空数组 = 清空截图）；未提供该字段则不修改
	if payload.Imgs != nil {
		game.Imgs = make([]in.Img, 0, len(payload.Imgs))
		for _, dataURLStr := range payload.Imgs {
			data, err := decodeDataURL(dataURLStr)
			if err != nil {
				return fmt.Errorf("解析游戏截图失败: %w", err)
			}
			game.Imgs = append(game.Imgs, data)
		}
	}
	return a.manager.EditGame(game, game)
}

func dataURL(data []byte) string {
	return "data:" + http.DetectContentType(data) + ";base64," + base64.StdEncoding.EncodeToString(data)
}

func decodeDataURL(s string) ([]byte, error) {
	if idx := strings.Index(s, "base64,"); idx >= 0 {
		s = s[idx+len("base64,"):]
	}
	return base64.StdEncoding.DecodeString(s)
}
