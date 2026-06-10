package main

import (
	"os"

	"android-tool-mvc/app/controllers"
	"android-tool-mvc/app/models"
	"android-tool-mvc/resources/views"
	"android-tool-mvc/routes"
)

func main() {
	// Initialize Models
	executor := models.NewADBExecutor()
	deviceModel := models.NewDeviceModel(executor)
	appModel := models.NewApplicationModel(executor)
	fileModel := models.NewFileTransferModel(executor)
	firmwareModel := models.NewFirmwareExtractor()

	// Initialize Views
	view := views.NewConsoleView()

	// Initialize Controllers
	deviceController := controllers.NewDeviceController(deviceModel, view)
	appController := controllers.NewAppController(appModel, view)
	fileController := controllers.NewFileController(fileModel, view)
	firmwareController := controllers.NewFirmwareController(firmwareModel, view)

	// Initialize Router
	router := routes.NewRouter(deviceController, appController, fileController, firmwareController, view, executor)

	// Route the command line arguments
	router.Route(os.Args)
}
