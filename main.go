package main

import (
	"context"
	"embed"
	runtime2 "github.com/wailsapp/wails/v2/pkg/runtime"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"devbox/backend/consts"
	"devbox/backend/service"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var version = "0.0.0"

func main() {
	// Create an instance of the service structure
	sysSvc := service.System()
	formatJSONSvc := service.NewFormatJSONService()
	prefSvc := service.Preferences()
	prefSvc.SetAppVersion(version)
	windowWidth, windowHeight, maximised := prefSvc.GetWindowSize()

	windowStartState := options.Normal
	if maximised {
		windowStartState = options.Maximised
	}

	// Menu
	appMenu := menu.NewMenu()
	if runtime.GOOS == "darwin" {
		appMenu.Append(menu.AppMenu())
		appMenu.Append(menu.EditMenu())
		appMenu.Append(menu.WindowMenu())
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:                    "Dev Box",
		Width:                    windowWidth,
		Height:                   windowHeight,
		MinWidth:                 consts.MinWindowWidth,
		MinHeight:                consts.MinWindowHeight,
		WindowStartState:         windowStartState,
		Frameless:                runtime.GOOS != "darwin",
		Menu:                     appMenu,
		EnableDefaultContextMenu: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			sysSvc.Start(ctx, version)
			formatJSONSvc.Startup(ctx)
			prefSvc.Startup(ctx)
		},
		OnDomReady: func(ctx context.Context) {
			x, y := prefSvc.GetWindowPosition(ctx)
			runtime2.WindowSetPosition(ctx, x, y)
			runtime2.WindowShow(ctx)
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			x, y := runtime2.WindowGetPosition(ctx)
			prefSvc.SaveWindowPosition(x, y)
			return false
		},
		OnShutdown: func(ctx context.Context) {
		},
		Bind: []interface{}{
			sysSvc,
			formatJSONSvc,
			prefSvc,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
			About: &mac.AboutInfo{
				Title:   "Dev Box " + version,
				Message: "A modern lightweight cross-platform Dev Box.\n\nCopyright © 2024",
				Icon:    icon,
			},
			WebviewIsTransparent: false,
			WindowIsTranslucent:  true,
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               true,
			DisableFramelessWindowDecorations: true,
		},
		Linux: &linux.Options{
			Icon:                icon,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyOnDemand,
			WindowIsTranslucent: true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
