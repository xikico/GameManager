import type {
  AddGameResult,
  CategoryDTO,
  EditGamePayload,
  GameDTO,
  SearchGameConditionDTO,
  SettingsDTO,
} from './types'

const app = () => window.go.wails_support.App

export const api = {
  getAllCategory: (): Promise<CategoryDTO[]> => app().GetAllCategory(),

  getGames: (condition: SearchGameConditionDTO): Promise<GameDTO[]> =>
    app().GetGames(condition),

  getGameByCategory: (categoryId: string): Promise<GameDTO[]> =>
    app().GetGameByCategory({ Id: categoryId, Name: '', Num: 0 }),

  getGameImgs: (gameId: string): Promise<string[]> => app().GetGameImgs(gameId),

  iconToBase64: (iconPath: string): Promise<string> => app().IconToBase64(iconPath),

  openGame: (game: GameDTO): Promise<void> => app().OpenGame(game),

  deleteGame: (game: GameDTO, delAll: boolean): Promise<void> =>
    app().DeleteGame(game, delAll),

  openFolder: (path: string): Promise<void> => app().OpenFolder(path),

  addGame: (path: string, password: string): Promise<AddGameResult> =>
    app().AddGame(path, password),

  confirmAddMany: (path: string): Promise<void> => app().ConfirmAddMany(path),

  editGame: (payload: EditGamePayload): Promise<void> => app().EditGame(payload),

  getSettings: (): Promise<SettingsDTO> => app().GetSettings(),

  updateSettings: (settings: SettingsDTO): Promise<void> => app().UpdateSettings(settings),
}

export function formatTime(t: string): string {
  if (!t) return ''
  const d = new Date(t)
  if (isNaN(d.getTime())) return ''
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}
