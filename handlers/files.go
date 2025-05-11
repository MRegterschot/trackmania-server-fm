package handlers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/MRegterschot/trackmania-server-fm/config"
	"github.com/MRegterschot/trackmania-server-fm/structs"
	"github.com/MRegterschot/trackmania-server-fm/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Handle files upload
func HandleUploadFiles(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("No files found")
	}

	paths := form.Value["paths[]"]
	if len(files) != len(paths) {
		return c.Status(fiber.StatusBadRequest).SendString("Number of files and paths do not match")
	}

	var errors []string

	for i, file := range files {
		relativePath := paths[i]

		// Set the destination path for the uploaded file
		dest := filepath.Join(config.AppEnv.UserDataPath, filepath.Clean("/"+relativePath))

		// Check if the path is in the UserData directory
		if !strings.HasPrefix(dest, config.AppEnv.UserDataPath) {
			errors = append(errors, "Invalid file path: "+file.Filename)
			continue
		}

		// Create the directory if it doesn't exist
		dir := filepath.Dir(dest)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			zap.L().Error("Error creating directory", zap.String("path", dest), zap.Error(err))
			errors = append(errors, "Failed to create directory: "+dest)
			continue
		}

		// Save the file to the destination
		if err := c.SaveFile(file, dest); err != nil {
			zap.L().Error("Error saving file", zap.String("path", dest), zap.Error(err))
			errors = append(errors, "Failed to save file: "+file.Filename)
			continue
		}

		zap.L().Info("File uploaded", zap.String("path", dest))
	}

	if len(errors) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Some files could not be uploaded",
			"errors":  errors,
		})
	}

	return c.SendString("Files uploaded successfully")
}

// Handle file and directory deletion
func HandleDeleteFiles(c *fiber.Ctx) error {
	// Get the file or directory paths from the request
	var paths []string
	if err := c.BodyParser(&paths); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	if len(paths) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("No file paths provided")
	}

	var errors []string

	for _, path := range paths {
		// Set the destination path for the uploaded file
		cleanPath := filepath.Join(config.AppEnv.UserDataPath, filepath.Clean("/"+path))

		// Check if the path is in the UserData directory
		if !strings.HasPrefix(cleanPath, config.AppEnv.UserDataPath) {
			errors = append(errors, "Invalid file path: "+path)
			continue
		}

		// Check if the file or directory exists before trying to delete it
		if _, err := os.Stat(cleanPath); err != nil {
			if os.IsNotExist(err) {
				errors = append(errors, "File/Directory does not exist: "+path)
				continue
			}
			zap.L().Error("Error checking file existence", zap.Error(err))
			errors = append(errors, "Error checking file existence: "+path)
			continue
		}

		// Delete the file or directory
		if err := os.RemoveAll(cleanPath); err != nil {
			zap.L().Error("Error deleting file or directory", zap.String("path", cleanPath), zap.Error(err))
			errors = append(errors, "Failed to delete: "+path)
			continue
		}

		zap.L().Info("File/Directory deleted", zap.String("path", cleanPath))
	}

	// If there are any errors, return them
	if len(errors) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Some files/directories could not be deleted",
			"errors":  errors,
		})
	}

	// Success message
	return c.SendString("Files/Directories deleted successfully")
}

// Handle file listing
func HandleListFiles(c *fiber.Ctx) error {
	relativePath := c.Params("*")
	absPath := filepath.Join(config.AppEnv.UserDataPath, filepath.Clean("/"+relativePath))

	// Prevent path traversal
	if !strings.HasPrefix(absPath, config.AppEnv.UserDataPath) {
		return c.Status(fiber.StatusForbidden).SendString("Invalid path")
	}

	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).SendString("Path not found")
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error accessing path")
	}

	// If it's a file, serve the file
	if !info.IsDir() {
		return c.SendFile(absPath)
	}

	// If it's a directory, return JSON listing
	entries, err := os.ReadDir(absPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read directory")
	}

	result := make([]structs.FileEntry, 0, len(entries))
	for _, entry := range entries {
		entryInfo, _ := entry.Info()
		result = append(result, structs.FileEntry{
			Name:         entry.Name(),
			Path:         filepath.Join("/UserData", relativePath, entry.Name()),
			IsDir:        entry.IsDir(),
			Size:         utils.GetSizeIfFile(entryInfo),
			LastModified: entryInfo.ModTime().UTC(),
		})
	}

	return c.JSON(result)
}
