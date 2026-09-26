package service

import (
	"GameManager/domain/entity"
	"GameManager/domain/ports/in"
	"GameManager/domain/ports/out/db"
	"GameManager/domain/ports/out/unzip"
	"GameManager/domain/ports/out/utils"
	domainUtils "GameManager/domain/utils"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type GameManager struct {
	db    db.DB
	utils utils.Utils
	unzip unzip.UnzipFunc
}

func NewGameManager(db db.DB, utils utils.Utils, unzip unzip.UnzipFunc) *GameManager {
	g := GameManager{
		db:    db,
		utils: utils,
		unzip: unzip,
	}
	return &g
}

// PredictGame 必须传入单个文件夹路径
func (g *GameManager) PredictGame(game in.GameDTO) (in.GameDTO, error) {
	normalizeGameDTOPaths(&game)
	gamePath := game.Path
	f, err := os.Open(gamePath)
	if err != nil {
		return game, fmt.Errorf("打开文件夹失败:%e", err)
	}
	defer f.Close()
	if stat, _ := f.Stat(); !stat.IsDir() {
		return game, errors.New("必须传入单个文件夹路径")
	}
	fileList, err := f.ReadDir(-1)
	if err != nil {
		return game, fmt.Errorf("遍历目录失败:%e", err)
	}

	exeList := make([]string, 0, 1)
	for _, childFile := range fileList {
		if !childFile.IsDir() {
			ext := path.Ext(childFile.Name())
			if strings.ToLower(ext) == ".exe" {
				exeList = append(exeList, filepath.Join(gamePath, childFile.Name()))
			}
		}
	}
	game.StartPath = predictExe(filepath.Base(gamePath), exeList)
	game.Name = filepath.Base(gamePath)
	return game, nil
}

func (g *GameManager) UnzipGame(path string, password string) (string, error) {
	return g.unzip(domainUtils.NormalizePath(path), password)
}

func (g *GameManager) SaveGame(games []in.GameDTO) error {
	if len(games) == 1 {
		gameDTO := games[0]
		normalizeGameDTOPaths(&gameDTO)
		return g.db.SaveGame(*domainUtils.GameToDomain(gameDTO))
	} else {
		gameDTOs := make([]entity.Game, len(games))
		for i, dto := range games {
			normalizeGameDTOPaths(&dto)
			gameDTOs[i] = *domainUtils.GameToDomain(dto)
		}
		return g.db.SaveGames(gameDTOs)
	}
}

func (g *GameManager) OpenGame(game in.GameDTO) error {
	normalizeGameDTOPaths(&game)
	return g.utils.OpenGame(domainUtils.GameToDomain(game))
}

func (g *GameManager) DeleteGame(game in.GameDTO, delAll bool) error {
	normalizeGameDTOPaths(&game)
	oldGame, newGame := domainUtils.GameToDomain(game), domainUtils.GameToDomain(game)
	newGame.IsDel = true
	err := g.db.EditGame(*oldGame, *newGame)
	if err != nil {
		return fmt.Errorf("数据库删除失败:%e", err)
	}
	if delAll {
		return g.utils.DeleteGame(oldGame)
	}
	return nil
}

func (g *GameManager) GetGameByCategory(category in.CategoryDTO) ([]in.GameDTO, error) {
	games, err := g.db.GetGameByCategory(*domainUtils.CategoryToDomain(category))
	if err != nil {
		return nil, fmt.Errorf("数据库查询失败:%e", err)
	}
	gameDTOs := g.sortAndConvertToDTO(games)
	return gameDTOs, nil
}

func (g *GameManager) GetGameImgs(gameId string) ([]in.Img, error) {
	imgs, err := g.db.GetGameImgs(gameId)
	if err != nil {
		return nil, fmt.Errorf("数据库查询失败:%e", err)
	}
	inImgs := make([]in.Img, len(imgs))
	for i, img := range imgs {
		inImgs[i] = in.Img(img)
	}
	return inImgs, nil
}

func (g *GameManager) GetGameByCondition(condition in.SearchGameConditionDTO) ([]in.GameDTO, error) {
	games, err := g.db.GetGameByCondition(*domainUtils.SearchGameConditionToDomain(condition))
	if err != nil {
		return nil, fmt.Errorf("数据库查询失败:%e", err)
	}
	gameDTOs := g.sortAndConvertToDTO(games)
	return gameDTOs, nil
}

func (g *GameManager) GetAllCategory() ([]in.CategoryDTO, error) {
	category, err := g.db.GetAllCategory()
	if err != nil {
		return nil, fmt.Errorf("数据库查询失败:%e", err)
	}
	categoryDTO := make([]in.CategoryDTO, len(category))
	for i, c := range category {
		categoryDTO[i] = domainUtils.CategoryToDTO(c)
	}
	return categoryDTO, nil
}

func (g *GameManager) EditGame(oldGame in.GameDTO, newGame in.GameDTO) error {
	normalizeGameDTOPaths(&oldGame)
	normalizeGameDTOPaths(&newGame)
	return g.db.EditGame(*domainUtils.GameToDomain(oldGame), *domainUtils.GameToDomain(newGame))
}

func (g *GameManager) GetSettings() (in.SettingsDTO, error) {
	settings, err := g.db.GetSettings()
	if err != nil {
		return in.SettingsDTO{}, fmt.Errorf("读取设置失败:%w", err)
	}
	g.utils.SetClipboardImageDetectionEnabled(settings.ClipboardImageDetectionEnabled)
	return in.SettingsDTO{ClipboardImageDetectionEnabled: settings.ClipboardImageDetectionEnabled}, nil
}

func (g *GameManager) UpdateSettings(settings in.SettingsDTO) error {
	domainSettings := entity.Settings{ClipboardImageDetectionEnabled: settings.ClipboardImageDetectionEnabled}
	if err := g.db.SaveSettings(domainSettings); err != nil {
		return fmt.Errorf("保存设置失败:%w", err)
	}
	g.utils.SetClipboardImageDetectionEnabled(settings.ClipboardImageDetectionEnabled)
	return nil
}

func normalizeGameDTOPaths(game *in.GameDTO) {
	game.IconPath = domainUtils.NormalizePath(game.IconPath)
	game.Path = domainUtils.NormalizePath(game.Path)
	game.StartPath = domainUtils.NormalizePath(game.StartPath)
}

// sortAndConvertToDTO 对游戏按插入时间降序排序并转换为DTO
func (g *GameManager) sortAndConvertToDTO(games []entity.Game) []in.GameDTO {
	// 按插入时间降序排序，最新的在前
	sort.Slice(games, func(i, j int) bool {
		return games[i].InsertTime.After(games[j].InsertTime)
	})
	gameDTOs := make([]in.GameDTO, len(games))
	for i, game := range games {
		gameDTOs[i] = domainUtils.GameToDTO(game)
	}
	return gameDTOs
}

var sensitiveRe = regexp.MustCompile(`(?i)(start|启动|开始|_chs)`)

func predictExe(fatherPath string, exeList []string) string {
	switch len(exeList) {
	case 1:
		return exeList[0]
	case 0:
		return ""
	}

	resStr := make([]string, 0, 1)

	for _, f := range exeList {
		childPath := filepath.Base(f)
		if strings.Contains(strings.ToLower(childPath), "unity") {
			continue
		}
		if strings.Contains(
			strings.ToLower(fatherPath),
			strings.ToLower(childPath),
		) || sensitiveRe.MatchString(childPath) {
			return f
		}
		resStr = append(resStr, f)
	}
	if len(resStr) == 0 {
		return exeList[0]
	}
	return resStr[0]
}
