package main

import (
	"PushMe/constant"
	"PushMe/internal/api"
	"PushMe/internal/services"
	"PushMe/internal/setting"
	"embed"
	_ "embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	if setting.Setting.Api.Enable {
		go api.Start()
	}
}

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        constant.AppName,
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(services.Notifier),
			application.NewService(&services.SettingService{}),
			application.NewService(&services.AppService{}),
			application.NewService(&services.MessageService{}),
			application.NewService(&services.PluginService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	windowShowing := false
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: constant.AppName,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(255, 255, 255),
		URL:              "/",
		Width:            constant.WindowWidth,
		Height:           constant.WindowHeight,
		Hidden:           true,
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
	})
	// window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
	// 	windowShowing = false
	// 	window.Hide()
	// 	e.Cancel()
	// })
	window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		windowShowing = false
		window.Hide()
		e.Cancel()
	})

	systemTray := app.SystemTray.New()
	menu := app.NewMenu()
	menu.Add("Quit").OnClick(func(data *application.Context) {
		app.Quit()
	})
	systemTray.SetMenu(menu)
	systemTray.SetTooltip(constant.AppName)

	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
	}

	systemTray.OnClick(func() {
		if windowShowing {
			window.Hide()
		} else {
			window.Show()
		}
		windowShowing = !windowShowing
	})

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
