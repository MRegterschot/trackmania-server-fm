package handlers

import (
	"regexp"

	"github.com/MRegterschot/trackmania-server-fm/config"
	"github.com/MRegterschot/trackmania-server-fm/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var mapsRe = regexp.MustCompile(`(?i)\.(Map|Challenge).*\.Gbx$`)

func HandleListMaps(c *fiber.Ctx) error {
	mapsDir := config.AppEnv.UserDataPath + "/Maps"

	mapFiles, err := utils.GetFilesRecursively(mapsDir)
	if err != nil {
		zap.L().Error("Error getting map files", zap.String("path", mapsDir), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get map files")
	}

	// Filter out files that are not .Map.Gbx or .Challenge.Gbx files
	filtered := make([]string, 0, len(mapFiles))
	for _, file := range mapFiles {
		if mapsRe.MatchString(file) {
			filtered = append(filtered, file)
		}
	}

	return c.JSON(filtered)
}
