package handlers

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/MRegterschot/trackmania-server-fm/config"
	"github.com/MRegterschot/trackmania-server-fm/structs"
	"github.com/gofiber/fiber/v2"
)

func download(t *testing.T, query string) (int, string, []byte) {
	t.Helper()
	app := fiber.New()
	app.Get("/download", HandleDownload)
	res, err := app.Test(httptest.NewRequest("GET", "/download?"+query, nil))
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(res.Body)
	return res.StatusCode, res.Header.Get("Content-Disposition"), body
}

func TestHandleDownload(t *testing.T) {
	root := t.TempDir()
	config.AppEnv = &structs.Env{UserDataPath: root}
	os.MkdirAll(filepath.Join(root, "Maps", "Sub"), 0o755)
	os.WriteFile(filepath.Join(root, "Maps", "a.txt"), []byte("aaa"), 0o644)
	os.WriteFile(filepath.Join(root, "Maps", "Sub", "b.txt"), []byte("bbb"), 0o644)

	status, disposition, body := download(t, "path=/UserData/Maps/a.txt")
	if status != 200 || string(body) != "aaa" || disposition == "" {
		t.Fatalf("file: status %d, disposition %q, body %q", status, disposition, body)
	}

	status, _, body = download(t, "path=/UserData/Maps")
	if status != 200 {
		t.Fatalf("folder: status %d", status)
	}
	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]bool{}
	for _, f := range zr.File {
		names[f.Name] = true
	}
	for _, want := range []string{"Maps/a.txt", "Maps/Sub/b.txt"} {
		if !names[want] {
			t.Errorf("zip is missing %s, has %v", want, names)
		}
	}

	if status, _, _ = download(t, "path=/UserData/nope"); status != 404 {
		t.Errorf("missing: status %d", status)
	}
	if status, _, _ = download(t, ""); status != 400 {
		t.Errorf("empty: status %d", status)
	}
	if status, _, _ = download(t, "path=/UserData/Maps/a.txt&path=/UserData/Maps/Sub/b.txt"); status != 200 {
		t.Errorf("several: status %d", status)
	}
}
