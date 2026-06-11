package controllers

import (
	"android-tool-mvc/app/models"
	"android-tool-mvc/resources/views"
)

type DiagnosticsController struct {
	model *models.DiagnosticsModel
	view  *views.ConsoleView
}

func NewDiagnosticsController(model *models.DiagnosticsModel, view *views.ConsoleView) *DiagnosticsController {
	return &DiagnosticsController{
		model: model,
		view:  view,
	}
}

func (c *DiagnosticsController) ShowAll(deviceID string) {
	report, err := c.model.GetAll(deviceID)
	if err != nil {
		c.view.RenderError(err)
		return
	}
	c.view.RenderDiagnostics(report)
}

func (c *DiagnosticsController) ShowMemory(deviceID string) {
	info, err := c.model.GetMemory(deviceID)
	if err != nil {
		c.view.RenderError(err)
		return
	}
	c.view.RenderDiagnostics(info)
}

func (c *DiagnosticsController) ShowCPU(deviceID string) {
	info, err := c.model.GetCPU(deviceID)
	if err != nil {
		c.view.RenderError(err)
		return
	}
	c.view.RenderDiagnostics(info)
}

func (c *DiagnosticsController) ShowStorage(deviceID string) {
	info, err := c.model.GetStorage(deviceID)
	if err != nil {
		c.view.RenderError(err)
		return
	}
	c.view.RenderDiagnostics(info)
}

func (c *DiagnosticsController) ShowDisplay(deviceID string) {
	info, err := c.model.GetDisplay(deviceID)
	if err != nil {
		c.view.RenderError(err)
		return
	}
	c.view.RenderDiagnostics(info)
}

func (c *DiagnosticsController) ShowNetwork(deviceID string) {
	info, err := c.model.GetNetwork(deviceID)
	if err != nil {
		c.view.RenderError(err)
		return
	}
	c.view.RenderDiagnostics(info)
}

func (c *DiagnosticsController) ShowSensors(deviceID string) {
	info, err := c.model.GetSensors(deviceID)
	if err != nil {
		c.view.RenderError(err)
		return
	}
	c.view.RenderDiagnostics(info)
}
