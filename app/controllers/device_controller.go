package controllers

import (
	"strings"

	"android-tool-mvc/app/models"
	"android-tool-mvc/resources/views"
)

// DeviceController handles device-related commands.
type DeviceController struct {
	model *models.DeviceModel
	view  *views.ConsoleView
}

// NewDeviceController creates a new DeviceController.
func NewDeviceController(model *models.DeviceModel, view *views.ConsoleView) *DeviceController {
	return &DeviceController{
		model: model,
		view:  view,
	}
}

// Index handles the "devices" command.
func (c *DeviceController) Index() {
	devices, err := c.model.All()
	if err != nil {
		c.view.RenderError(err)
		return
	}
	c.view.RenderDevicesTable(devices)
}

// Reboot handles the "reboot" command.
func (c *DeviceController) Reboot(deviceID string) {
	_, err := c.model.Reboot(deviceID)
	if err != nil {
		c.view.RenderError(err)
		return
	}
	
	msg := "Device reboot triggered"
	if deviceID != "" {
		msg += " for device " + deviceID
	}
	c.view.RenderSuccess(msg)
}

// ShowDeviceInfo displays detailed information for a device.
func (c *DeviceController) ShowDeviceInfo(deviceID string) {
	info, err := c.model.GetDetailedInfo(deviceID)
	if err != nil {
		c.view.RenderError(err)
		return
	}
	c.view.RenderDeviceInfo(info)
}

// Screenshot captures the device screen and saves to a PNG file.
func (c *DeviceController) Screenshot(deviceID string, outputPath string) {
	if strings.TrimSpace(outputPath) == "" {
		outputPath = ""
	}
	path, err := c.model.Screenshot(deviceID, outputPath)
	if err != nil {
		c.view.RenderError(err)
		return
	}
	c.view.RenderSuccess("Screenshot saved to " + path)
}

