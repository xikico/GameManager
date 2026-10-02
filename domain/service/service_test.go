package service

import (
	"GameManager/adapters/out/db/fake"
	"GameManager/domain/entity"
	"GameManager/domain/ports/in"
	"GameManager/domain/ports/out/filesystem"
	"errors"
	"testing"
)

type gameOperatorStub struct {
	opened  *entity.Game
	deleted *entity.Game
}

func (s *gameOperatorStub) OpenGame(game *entity.Game) error {
	s.opened = game
	return nil
}

func (s *gameOperatorStub) DeleteGame(game *entity.Game) error {
	s.deleted = game
	return nil
}

func (s *gameOperatorStub) OpenFolder(path string) error { return nil }

type clipboardSettingsStub struct{ enabled bool }

func (s *clipboardSettingsStub) SetClipboardImageDetectionEnabled(enabled bool) { s.enabled = enabled }

type inspectorStub struct {
	snapshots map[string]filesystem.Snapshot
}

func (s inspectorStub) Inspect(path string) (filesystem.Snapshot, error) {
	snapshot, ok := s.snapshots[path]
	if !ok {
		return filesystem.Snapshot{}, errors.New("path not found")
	}
	return snapshot, nil
}

func testManager(inspector filesystem.Inspector) (*GameManager, *gameOperatorStub) {
	games := &gameOperatorStub{}
	clipboard := &clipboardSettingsStub{}
	unzip := func(path, password string) (string, error) { return "", errors.New("unexpected unzip") }
	return NewGameManager(fake.GetDB(), games, games, clipboard, unzip, inspector), games
}

func TestPredictGameUsesDirectorySnapshot(t *testing.T) {
	manager, _ := testManager(inspectorStub{snapshots: map[string]filesystem.Snapshot{
		`C:\Games\Demo`: {
			Path: `C:\Games\Demo`, Base: "Demo", IsDir: true,
			Entries: []filesystem.Entry{{Name: "Demo.exe"}, {Name: "readme.txt"}},
		},
	}})

	game, err := manager.PredictGame(in.GameDTO{Path: ` "C:\Games\Demo" `})
	if err != nil {
		t.Fatalf("PredictGame: %v", err)
	}
	if game.Name != "Demo" || game.StartPath != `C:\Games\Demo\Demo.exe` {
		t.Fatalf("unexpected prediction: %#v", game)
	}
}

func TestImportGameDetectsBatch(t *testing.T) {
	manager, _ := testManager(inspectorStub{snapshots: map[string]filesystem.Snapshot{
		`C:\Collection`: {
			Path: `C:\Collection`, Base: "Collection", IsDir: true,
			Entries: []filesystem.Entry{{Name: "One", IsDir: true}, {Name: "Two", IsDir: true}},
		},
	}})

	result, err := manager.ImportGame(`C:\Collection`, "")
	if err != nil {
		t.Fatalf("ImportGame: %v", err)
	}
	if !result.IsMany {
		t.Fatalf("expected batch result: %#v", result)
	}
}

func TestOpenGameNormalizesPaths(t *testing.T) {
	manager, games := testManager(inspectorStub{})
	if err := manager.OpenGame(in.GameDTO{Path: ` "C:\Game" `, StartPath: ` "C:\Game\game.exe" `}); err != nil {
		t.Fatalf("OpenGame: %v", err)
	}
	if games.opened.Path != `C:\Game` || games.opened.StartPath != `C:\Game\game.exe` {
		t.Fatalf("paths were not normalized: %#v", games.opened)
	}
}

func TestSaveGameRejectsEmptyInput(t *testing.T) {
	manager, _ := testManager(inspectorStub{})
	if err := manager.SaveGame(nil); err == nil {
		t.Fatal("expected empty input error")
	}
}
