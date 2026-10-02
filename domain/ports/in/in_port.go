package in

type GameManager interface {
	PredictGame(game GameDTO) (GameDTO, error) // 自动推断导入游戏信息
	ImportGame(path string, password string) (ImportGameResult, error)
	ConfirmImportMany(path string) error
	UnzipGame(path string, password string) (string, error)
	SaveGame(games []GameDTO) error
	OpenGame(game GameDTO) error
	OpenFolder(path string) error
	DeleteGame(game GameDTO, delAll bool) error
	GetGameByCategory(category CategoryDTO) ([]GameDTO, error)
	GetGameImgs(gameId string) ([]Img, error) // 除了这个之外的GetGame都不返回图片（imgs）
	GetGameByCondition(condition SearchGameConditionDTO) ([]GameDTO, error)
	GetAllCategory() ([]CategoryDTO, error)
	EditGame(oldGame GameDTO, newGame GameDTO) error
	GetSettings() (SettingsDTO, error)
	InitializeSettings() error
	UpdateSettings(settings SettingsDTO) error
}
