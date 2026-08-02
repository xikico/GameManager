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
  NewImgs: string[]
}

export interface SeriesGroup {
  name: string
  games: GameDTO[]
  latest: GameDTO
}
