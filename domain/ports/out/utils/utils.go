package utils

import (
	"GameManager/domain/entity"
	"context"
)

type Utils interface {
	OpenGame(game *entity.Game) error
	DeleteGame(game *entity.Game) error
	Img2Base64(ctx context.Context)
}
