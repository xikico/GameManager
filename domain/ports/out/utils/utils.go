package utils

import (
	"GameManager/domain/entity"
	"context"
)

type GameOperator interface {
	OpenGame(game *entity.Game) error
	DeleteGame(game *entity.Game) error
}

type FolderOpener interface {
	OpenFolder(path string) error
}

type ClipboardSettings interface {
	SetClipboardImageDetectionEnabled(enabled bool)
}

type ClipboardMonitor interface {
	Img2Base64(ctx context.Context)
}
