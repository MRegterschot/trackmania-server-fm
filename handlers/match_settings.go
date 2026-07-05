package handlers

import (
	"regexp"

	"github.com/MRegterschot/trackmania-server-fm/config"
	"github.com/MRegterschot/trackmania-server-fm/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var matchSettingsRe = regexp.MustCompile(`\.txt$`)

func HandleListMatchSettings(c *fiber.Ctx) error {
	matchSettingsDir := config.AppEnv.UserDataPath + "/Maps/MatchSettings"

	matchSettingsFiles, err := utils.GetFilesRecursively(matchSettingsDir)
	if err != nil {
		zap.L().Error("Error getting match settings files", zap.String("path", matchSettingsDir), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get match settings files")
	}

	// Filter out files that are not .txt files
	filtered := make([]string, 0, len(matchSettingsFiles))
	for _, file := range matchSettingsFiles {
		if matchSettingsRe.MatchString(file) {
			filtered = append(filtered, file)
		}
	}

	return c.JSON(filtered)
}
