package main

import (
	"TotalControl/backend/scripting"
	"TotalControl/backend/utils"
	"embed"
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var pluginManager *scripting.PluginManager

func GetPluginManager() *scripting.PluginManager {
	return pluginManager
}

func ParseWailsLogLevel(level string) logger.LogLevel {
	switch level {
	case "debug":
		return logger.DEBUG
	case "info":
		return logger.INFO
	case "warn":
		return logger.WARNING
	case "error":
		return logger.ERROR
	default:
		log.Warnf("Invalid Wails log level '%s', defaulting to 'info'", level)
		return logger.INFO
	}
}

func main() {
	logPath := flag.String("log-path", "", "Path to log file (default: stdout)")
	logLevel := flag.String("log-level", "debug", "Log level (debug, info, warn, error, fatal, panic)")
	appDataPath := flag.String("app-data-path", "./", "Path to application data directory (default: system default path)")
	wailsLogLevelStr := flag.String("wails-log-level", "info", "Wails log level (debug, info, warn, error, fatal, panic)")
	flag.Parse()

	if *logPath != "" {
		f, err := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		log.SetOutput(f)
	}

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Warnf("Invalid log level '%s', defaulting to 'info'", *logLevel)
		level = log.InfoLevel
	}
	log.SetLevel(level)
	log.SetReportCaller(true)
	log.SetFormatter(&utils.CustomFormatter{})

	// Initialize the plugin manager
	pluginDir := filepath.Join(utils.GetAppDataPath(), "plugins")
	if *appDataPath != "" {
		pluginDir = filepath.Join(*appDataPath, "plugins")
	}
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		if err := os.MkdirAll(pluginDir, 0755); err != nil {
			log.Fatalf("Failed to create plugin directory: %v", err)
		}
	}
	pluginManager, err = scripting.NewPluginManager(pluginDir)
	if err != nil {
		log.Fatalf("Failed to initialize plugin manager: %v", err)
	}

	// Create an instance of the app structure
	app := NewApp()
	AppMenu := createMenus(app)

	var wailsLogLevel logger.LogLevel
	if *wailsLogLevelStr != "" {
		logLevel := ParseWailsLogLevel(*wailsLogLevelStr)
		log.Infof("Setting Wails log level to %s", logLevel)
		wailsLogLevel = logLevel
	} else {
		log.Warn("No Wails log level specified, defaulting to 'info'")
		wailsLogLevel = logger.INFO
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:     "TotalControl",
		Width:     900,
		Height:    600,
		MinWidth:  900,
		MinHeight: 600,
		//MaxWidth:          1200,
		//MaxHeight:         800,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		//BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 0},
		Menu:          AppMenu,
		Logger:        nil,
		LogLevel:      wailsLogLevel,
		OnStartup:     app.startup,
		OnDomReady:    app.domReady,
		OnBeforeClose: app.beforeClose,
		OnShutdown:    app.shutdown,
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               "de95813c-cf49-446c-b34e-92fb592503e5",
			OnSecondInstanceLaunch: app.onSecondInstanceLaunch,
		},
		WindowStartState: options.Normal,
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Handler:    nil,
			Middleware: nil,
		},
		Bind: []interface{}{
			app,
		},
		EnumBind:       []interface{}{},
		ErrorFormatter: func(err error) any { return err.Error() },
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               false,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme:                             windows.SystemDefault,
		},
		// Mac platform specific options
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  true,
				HideTitleBar:               false,
				FullSizeContent:            true,
				UseToolbar:                 false,
				HideToolbarSeparator:       false,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
		Linux: &linux.Options{
			Icon: icon,
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: true,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

func createMenus(app *App) *menu.Menu {
	AppMenu := menu.NewMenu()
	if runtime.GOOS == "darwin" {
		AppMenu.Append(menu.AppMenu()) // On macOS platform, this must be done right after `NewMenu()`
	}
	if runtime.GOOS == "darwin" {
		AppMenu.Append(menu.EditMenu())
	}

	HelpMenu := AppMenu.AddSubmenu("Help")
	HelpMenu.AddText("About", keys.CmdOrCtrl("a"), func(_ *menu.CallbackData) {
		rt.BrowserOpenURL(app.ctx, "https://github.com/subtixx/total-control")
	})
	return AppMenu
}
