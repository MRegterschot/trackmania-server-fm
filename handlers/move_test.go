package handlers

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MRegterschot/trackmania-server-fm/config"
	"github.com/MRegterschot/trackmania-server-fm/structs"
	"github.com/gofiber/fiber/v2"
)

func moveRequest(t *testing.T, body string) int {
	t.Helper()
	app := fiber.New()
	app.Post("/move", HandleMoveItems)
	req := httptest.NewRequest("POST", "/move", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	return res.StatusCode
}

func TestHandleMoveItems(t *testing.T) {
	root := t.TempDir()
	config.AppEnv = &structs.Env{UserDataPath: root}
	os.MkdirAll(filepath.Join(root, "Maps", "Sub"), 0o755)
	os.WriteFile(filepath.Join(root, "Maps", "a.txt"), []byte("a"), 0o644)
	os.WriteFile(filepath.Join(root, "Maps", "b.txt"), []byte("b"), 0o644)

	if got := moveRequest(t, `{"items":[{"from":"/UserData/Maps/a.txt","to":"/UserData/Maps/Sub/c.txt"}]}`); got != 200 {
		t.Fatalf("move: status %d", got)
	}
	if _, err := os.Stat(filepath.Join(root, "Maps", "Sub", "c.txt")); err != nil {
		t.Fatal("file was not moved")
	}

	cases := map[string]struct {
		body   string
		status int
	}{
		"exists":    {`{"items":[{"from":"/UserData/Maps/b.txt","to":"/UserData/Maps/Sub/c.txt"}]}`, 409},
		"traversal": {`{"items":[{"from":"/UserData/Maps/b.txt","to":"/UserData/../x.txt"}]}`, 200},
		"root":      {`{"items":[{"from":"/UserData","to":"/UserData/x"}]}`, 403},
		"into self": {`{"items":[{"from":"/UserData/Maps","to":"/UserData/Maps/Sub/Maps"}]}`, 400},
		"missing":   {`{"items":[{"from":"/UserData/Maps/nope.txt","to":"/UserData/Maps/n.txt"}]}`, 404},
	}
	for name, c := range cases {
		if got := moveRequest(t, c.body); got != c.status {
			t.Errorf("%s: status %d, want %d", name, got, c.status)
		}
	}
	// "../x.txt" is cleaned into UserData rather than escaping it
	if _, err := os.Stat(filepath.Join(filepath.Dir(root), "x.txt")); err == nil {
		t.Error("moved outside UserData")
	}
}
