export interface CategoryDTO {
  Id: string
  Name: string
  Num: number
}

export interface GameDTO {
  Id: string
  IconPath: string
  Name: string
  NickName: string
  Series: string
  Description: string
  Path: string
  StartPath: string
  Category: CategoryDTO
  Imgs: string | null
  IsPlay: boolean
  IsDel: boolean
  InsertTime: string
}

export interface SearchGameConditionDTO {
  Name: string
  Series: string
  Description: string
  Id: string
  IsPlay: boolean | null
  InsertTimeStart: string | null
  InsertTimeEnd: string | null
}

export interface AddGameResult {
  is_many: boolean
  need_pass: boolean
  message: string
}

export interface SettingsDTO {
  ClipboardImageDetectionEnabled: boolean
}

export interface EditGamePayload {
  Id: string
  IconPath: string
  Name: string
  NickName: string
  Series: string
  Description: string
  Path: string
  StartPath: string
  CategoryId: string
  IsPlay: boolean
  // Imgs 完整的最终展示图列表（data URL），全量替换；传空数组则清空截图
  Imgs: string[]
}

export interface SeriesGroup {
  name: string
  games: GameDTO[]
  latest: GameDTO
}
