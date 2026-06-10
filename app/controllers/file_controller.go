package controllers

import (
	"errors"
	"strings"

	"android-tool-mvc/app/models"
	"android-tool-mvc/resources/views"
)

// FileController handles file operations like pushing and pulling files.
type FileController struct {
	model *models.FileTransferModel
	view  *views.ConsoleView
}

// NewFileController creates a new FileController.
func NewFileController(model *models.FileTransferModel, view *views.ConsoleView) *FileController {
	return &FileController{
		model: model,
		view:  view,
	}
}

// Push handles uploading a local file/directory to the device.
func (c *FileController) Push(localPath string, remotePath string, deviceID string) {
	if strings.TrimSpace(localPath) == "" || strings.TrimSpace(remotePath) == "" {
		c.view.RenderError(errors.New("missing arguments: push requires both <local-path> and <remote-path>"))
		return
	}

	_, err := c.model.Push(localPath, remotePath, deviceID)
	if err != nil {
		c.view.RenderError(err)
		return
	}

	msg := "Successfully pushed " + localPath + " to " + remotePath
	if deviceID != "" {
		msg += " on device " + deviceID
	}
	c.view.RenderSuccess(msg)
}

// Pull handles downloading a remote file/directory from the device.
func (c *FileController) Pull(remotePath string, localPath string, deviceID string) {
	if strings.TrimSpace(remotePath) == "" || strings.TrimSpace(localPath) == "" {
		c.view.RenderError(errors.New("missing arguments: pull requires both <remote-path> and <local-path>"))
		return
	}

	_, err := c.model.Pull(remotePath, localPath, deviceID)
	if err != nil {
		c.view.RenderError(err)
		return
	}

	msg := "Successfully pulled " + remotePath + " to " + localPath
	if deviceID != "" {
		msg += " from device " + deviceID
	}
	c.view.RenderSuccess(msg)
}
