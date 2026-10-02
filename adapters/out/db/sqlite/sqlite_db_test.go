package sqlite

import (
	"GameManager/domain/entity"
	"path/filepath"
	"testing"
)

func testDB(t *testing.T) *sqliteDB {
	t.Helper()
	dsn := "file:" + filepath.ToSlash(filepath.Join(t.TempDir(), "test.db")) + "?_busy_timeout=5000"
	db, err := openSqliteDB(dsn)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() {
		if err := db.close(); err != nil {
			t.Errorf("close test database: %v", err)
		}
	})
	return db
}

func TestSaveGamesRejectsEmptyInput(t *testing.T) {
	if err := testDB(t).SaveGames(nil); err == nil {
		t.Fatal("expected empty batch error")
	}
}

func TestEditGameCreatesCategoryAndMovesCount(t *testing.T) {
	db := testDB(t)
	game := entity.Game{Name: "Game", Path: `C:\Game`, Category: entity.Category{Id: "000000"}}
	if err := db.SaveGame(game); err != nil {
		t.Fatalf("save game: %v", err)
	}

	var saved struct{ Id string }
	if err := db.db.Table("game").Select("id").Where("name = ?", "Game").First(&saved).Error; err != nil {
		t.Fatalf("find game: %v", err)
	}
	game.Id = saved.Id
	updated := game
	updated.Category = entity.Category{Name: "New Category"}
	if err := db.EditGame(game, updated); err != nil {
		t.Fatalf("edit game: %v", err)
	}

	categories, err := db.GetAllCategory()
	if err != nil {
		t.Fatalf("get categories: %v", err)
	}
	counts := make(map[string]int, len(categories))
	for _, category := range categories {
		counts[category.Name] = category.Num
	}
	if counts["未分类"] != 0 || counts["New Category"] != 1 {
		t.Fatalf("unexpected category counts: %#v", counts)
	}
}

func TestRepeatedDeleteDoesNotDecrementTwice(t *testing.T) {
	db := testDB(t)
	game := entity.Game{Name: "Game", Path: `C:\Game`, Category: entity.Category{Id: "000000"}}
	if err := db.SaveGame(game); err != nil {
		t.Fatalf("save game: %v", err)
	}
	var saved struct{ Id string }
	if err := db.db.Table("game").Select("id").Where("name = ?", "Game").First(&saved).Error; err != nil {
		t.Fatalf("find game: %v", err)
	}
	game.Id = saved.Id
	deleted := game
	deleted.IsDel = true
	if err := db.EditGame(game, deleted); err != nil {
		t.Fatalf("first delete: %v", err)
	}
	if err := db.EditGame(game, deleted); err != nil {
		t.Fatalf("second delete: %v", err)
	}
	var category struct{ Num int }
	if err := db.db.Table("category").Select("num").Where("id = ?", "000000").First(&category).Error; err != nil {
		t.Fatalf("find category: %v", err)
	}
	if category.Num != 0 {
		t.Fatalf("category count = %d, want 0", category.Num)
	}
}
