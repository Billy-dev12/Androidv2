package controllers

import (
	"errors"
	"strings"

	"android-tool-mvc/app/models"
	"android-tool-mvc/resources/views"
)

// AppController handles app-related commands like install and uninstall.
type AppController struct {
	model *models.ApplicationModel
	view  *views.ConsoleView
}

// NewAppController creates a new AppController.
func NewAppController(model *models.ApplicationModel, view *views.ConsoleView) *AppController {
	return &AppController{
		model: model,
		view:  view,
	}
}

// Install handles the "install" command.
func (c *AppController) Install(apkPath string, deviceID string) {
	if strings.TrimSpace(apkPath) == "" {
		c.view.RenderError(errors.New("missing required argument: <apk-path>"))
		return
	}

	_, err := c.model.Install(apkPath, deviceID)
	if err != nil {
		c.view.RenderError(err)
		return
	}

	msg := "Successfully installed " + apkPath
	if deviceID != "" {
		msg += " on device " + deviceID
	}
	c.view.RenderSuccess(msg)
}

// Uninstall handles the "uninstall" command.
func (c *AppController) Uninstall(packageName string, deviceID string) {
	if strings.TrimSpace(packageName) == "" {
		c.view.RenderError(errors.New("missing required argument: <package-name>"))
		return
	}

	_, err := c.model.Uninstall(packageName, deviceID)
	if err != nil {
		c.view.RenderError(err)
		return
	}

	msg := "Successfully uninstalled " + packageName
	if deviceID != "" {
		msg += " from device " + deviceID
	}
	c.view.RenderSuccess(msg)
}
