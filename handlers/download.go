package handlers

import (
	"archive/zip"
	"bufio"
	"io"
	"io/fs"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/MRegterschot/trackmania-server-fm/config"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Resolves a UserData-relative path, allowing the UserData root itself
func resolveDownloadPath(relative string) (string, bool) {
	decoded, err := url.PathUnescape(relative)
	if err != nil {
		return "", false
	}

	root := filepath.Clean(config.AppEnv.UserDataPath)
	abs := filepath.Join(root, filepath.Clean("/"+strings.TrimPrefix(decoded, "/UserData")))
	if abs != root && !strings.HasPrefix(abs, root+string(filepath.Separator)) {
		return "", false
	}
	return abs, true
}

func addToZip(zw *zip.Writer, abs, name string) error {
	return filepath.WalkDir(abs, func(current string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(abs, current)
		if err != nil {
			return err
		}
		entryName := filepath.ToSlash(filepath.Join(name, rel))

		if d.IsDir() {
			_, err := zw.Create(entryName + "/")
			return err
		}
		if !d.Type().IsRegular() {
			return nil
		}

		w, err := zw.Create(entryName)
		if err != nil {
			return err
		}
		f, err := os.Open(current)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(w, f)
		return err
	})
}

// Download one file as is, or any folder / several items as a zip
func HandleDownload(c *fiber.Ctx) error {
	paths := c.Context().QueryArgs().PeekMulti("path")
	if len(paths) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("No paths provided")
	}

	type item struct {
		abs  string
		name string
		dir  bool
	}
	items := make([]item, 0, len(paths))
	seen := map[string]bool{}

	for _, raw := range paths {
		abs, ok := resolveDownloadPath(string(raw))
		if !ok {
			return c.Status(fiber.StatusForbidden).SendString("Invalid path")
		}

		info, err := os.Stat(abs)
		if os.IsNotExist(err) {
			return c.Status(fiber.StatusNotFound).SendString("Path not found")
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error accessing path")
		}

		name := filepath.Base(abs)
		if abs == filepath.Clean(config.AppEnv.UserDataPath) {
			name = "UserData"
		}
		if seen[name] {
			return c.Status(fiber.StatusBadRequest).SendString("Duplicate item name: " + name)
		}
		seen[name] = true
		items = append(items, item{abs, name, info.IsDir()})
	}

	if len(items) == 1 && !items[0].dir {
		c.Attachment(items[0].name)
		return c.SendFile(items[0].abs)
	}

	archiveName := items[0].name
	if len(items) > 1 {
		archiveName = "download"
	}

	c.Set(fiber.HeaderContentType, "application/zip")
	c.Set(fiber.HeaderContentDisposition, mime.FormatMediaType("attachment", map[string]string{"filename": archiveName + ".zip"}))
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		zw := zip.NewWriter(w)
		for _, it := range items {
			if err := addToZip(zw, it.abs, it.name); err != nil {
				zap.L().Error("Error zipping item", zap.String("path", it.abs), zap.Error(err))
				break
			}
		}
		if err := zw.Close(); err != nil {
			zap.L().Error("Error closing zip", zap.Error(err))
		}
		w.Flush()
	})
	return nil
}
