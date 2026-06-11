package controllers

import (
	"android-tool-mvc/app/models"
	"android-tool-mvc/resources/views"
)

type SystemController struct {
	envModel     *models.EnvironmentModel
	configModel  *models.ConfigModel
	historyModel *models.HistoryModel
	view         *views.ConsoleView
}

func NewSystemController(
	envModel *models.EnvironmentModel,
	configModel *models.ConfigModel,
	historyModel *models.HistoryModel,
	view *views.ConsoleView,
) *SystemController {
	return &SystemController{
		envModel:     envModel,
		configModel:  configModel,
		historyModel: historyModel,
		view:         view,
	}
}

func (c *SystemController) ShowEnv() {
	tools, sysInfo := c.envModel.CheckAll()
	report := sysInfo + "\n" + c.envModel.ToolInfoToString(tools)
	c.view.RenderSystemInfo(report)
}

func (c *SystemController) ShowConfig() {
	_ = c.configModel.Load()
	cfg := c.configModel.Show()
	c.view.RenderSystemInfo(cfg)
}

func (c *SystemController) SetConfig(key, value string) {
	_ = c.configModel.Load()
	err := c.configModel.Set(key, value)
	if err != nil {
		c.view.RenderError(err)
		return
	}
	c.view.RenderSuccess("Config updated: " + key + " = " + value)
}

func (c *SystemController) ShowHistory() {
	content, err := c.historyModel.Show()
	if err != nil {
		c.view.RenderError(err)
		return
	}
	c.view.RenderSystemInfo(content)
}

func (c *SystemController) ClearHistory() {
	err := c.historyModel.Clear()
	if err != nil {
		c.view.RenderError(err)
		return
	}
	c.view.RenderSuccess("History cleared.")
}
