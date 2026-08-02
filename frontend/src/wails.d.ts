import type {
  AddGameResult,
  CategoryDTO,
  EditGamePayload,
  GameDTO,
  SearchGameConditionDTO,
} from './types'

declare global {
  interface Window {
    go: {
      wails_support: {
        App: {
          GetAllCategory(): Promise<CategoryDTO[]>
          GetGames(condition: SearchGameConditionDTO): Promise<GameDTO[]>
          GetGameByCategory(category: CategoryDTO): Promise<GameDTO[]>
          GetGameImgs(gameId: string): Promise<string[]>
          IconToBase64(iconPath: string): Promise<string>
          OpenGame(game: GameDTO): Promise<void>
          DeleteGame(game: GameDTO, delAll: boolean): Promise<void>
          OpenFolder(path: string): Promise<void>
          AddGame(path: string, password: string): Promise<AddGameResult>
          ConfirmAddMany(path: string): Promise<void>
          EditGame(payload: EditGamePayload): Promise<void>
        }
      }
    }
  }
}

export {}
