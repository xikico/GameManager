// Package wails_support 提供桌面端入口适配器。
//
// 该适配器是六边形架构的 in 侧新端口实现：
// 只负责把 Wails 前端的方法调用翻译成 in.GameManager 接口调用，
// 自身不包含任何业务逻辑，核心业务全部复用 domain/service 层。
package wails_support

import (
	"GameManager/domain/ports/in"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	if path == "" {
		return errors.New("路径不能为空")
	}
	cmd := exec.Command("explorer", path)
	return cmd.Start()
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
	pathType, err := analysisPath(path)
	if err != nil {
		return AddGameResult{}, err
	}
	switch pathType {
	case singleGameFolder:
		if err := a.processSingleGame(path); err != nil {
			return AddGameResult{}, err
		}
		return AddGameResult{Message: "已加入游戏"}, nil
	case manyGameFolder:
		return AddGameResult{IsMany: true, Message: "是否将其作为一个分类进行批量添加"}, nil
	case zipFile:
		dir, err := a.manager.UnzipGame(path, password)
		if err != nil {
			if err.Error() == "密码错误或需要密码但未提供" {
				return AddGameResult{NeedPass: true, Message: "需要密码或密码错误"}, nil
			}
			return AddGameResult{}, err
		}
		if err := a.processSingleGame(dir); err != nil {
			return AddGameResult{}, err
		}
		return AddGameResult{Message: "已加入游戏"}, nil
	}
	return AddGameResult{}, errors.New("无法识别的路径类型")
}

// ConfirmAddMany 将文件夹内的所有子目录批量添加，并把文件夹名作为分类
func (a *App) ConfirmAddMany(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("打不开文件夹：%w", err)
	}
	defer f.Close()

	entries, err := f.ReadDir(-1)
	if err != nil {
		return fmt.Errorf("读取文件夹失败：%w", err)
	}

	games := make([]in.GameDTO, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		game := in.GameDTO{
			Path: filepath.Join(path, entry.Name()),
			Category: in.CategoryDTO{
				Name: filepath.Base(path),
			},
		}
		game, err = a.manager.PredictGame(game)
		if err != nil {
			continue
		}
		games = append(games, game)
	}
	if len(games) == 0 {
		return errors.New("文件夹内没有找到可添加的游戏")
	}
	return a.manager.SaveGame(games)
}

// PredictGame 预测游戏启动路径等信息（不保存）
func (a *App) PredictGame(game in.GameDTO) (in.GameDTO, error) {
	return a.manager.PredictGame(game)
}

// ---------- 编辑游戏 ----------

// EditGamePayload 编辑游戏的请求参数（与原 gin 表单接口等价，改用 JSON 传输）
type EditGamePayload struct {
	Id          string
	IconPath    string
	Name        string
	NickName    string
	Series      string
	Description string
	Path        string
	StartPath   string
	CategoryId  string
	IsPlay      bool
	// NewImgs 新增的展示图，data URL 形式（仅提交新增/替换的图片）
	NewImgs []string
}

func (a *App) EditGame(payload EditGamePayload) error {
	game := in.GameDTO{
		Id:          payload.Id,
		IconPath:    payload.IconPath,
		Name:        payload.Name,
		NickName:    payload.NickName,
		Series:      payload.Series,
		Description: payload.Description,
		Path:        strings.ReplaceAll(payload.Path, "\"", ""),
		StartPath:   strings.ReplaceAll(payload.StartPath, "\"", ""),
		Category:    in.CategoryDTO{Id: payload.CategoryId},
		IsPlay:      payload.IsPlay,
	}
	for _, dataURLStr := range payload.NewImgs {
		data, err := decodeDataURL(dataURLStr)
		if err == nil {
			game.Imgs = append(game.Imgs, data)
		}
	}
	return a.manager.EditGame(game, game)
}

// ---------- 内部工具 ----------

func (a *App) processSingleGame(path string) error {
	game, err := a.manager.PredictGame(in.GameDTO{Path: path})
	if err != nil {
		return err
	}
	return a.manager.SaveGame([]in.GameDTO{game})
}

var dataURLPrefix = "data:image/jpeg;base64,"

func dataURL(data []byte) string {
	return dataURLPrefix + base64.StdEncoding.EncodeToString(data)
}

func decodeDataURL(s string) ([]byte, error) {
	if idx := strings.Index(s, "base64,"); idx >= 0 {
		s = s[idx+len("base64,"):]
	}
	return base64.StdEncoding.DecodeString(s)
}

// ---------- 路径类型分析（与 gin 适配器共用同一逻辑） ----------

type pathType int

const (
	singleGameFolder pathType = iota
	manyGameFolder
	zipFile
)

func analysisPath(filePath string) (pathType, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return zipFile, fmt.Errorf("获取文件信息失败：%e", err)
	}
	if !info.IsDir() {
		return zipFile, nil
	}
	f, err := os.Open(filePath)
	if err != nil {
		return zipFile, fmt.Errorf("打开文件夹失败：%e", err)
	}
	defer f.Close()

	entries, err := f.ReadDir(-1)
	if err != nil {
		return zipFile, fmt.Errorf("读取所有目录失败：%e", err)
	}
	if len(entries) == 0 {
		return zipFile, errors.New("传入的是个空目录")
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			return singleGameFolder, nil
		}
	}
	return manyGameFolder, nil
}
