package handlers

import (
	"regexp"

	"github.com/MRegterschot/trackmania-server-fm/config"
	"github.com/MRegterschot/trackmania-server-fm/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var scriptsRe = regexp.MustCompile(`\.Script\.txt$`)

func HandleListScripts(c *fiber.Ctx) error {
	scriptsDir := config.AppEnv.UserDataPath + "/Scripts/Modes"

	scriptFiles, err := utils.GetFilesRecursively(scriptsDir)
	if err != nil {
		zap.L().Error("Error getting script files", zap.String("path", scriptsDir), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get script files")
	}

	// Filter out files that are not .Script.txt files
	filtered := make([]string, 0, len(scriptFiles))
	for _, file := range scriptFiles {
		if scriptsRe.MatchString(file) {
			filtered = append(filtered, file)
		}
	}

	return c.JSON(filtered)
}
