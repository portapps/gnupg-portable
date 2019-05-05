//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico
package main

import (
	"fmt"

	. "github.com/portapps/portapps"
	"github.com/portapps/portapps/pkg/dialog"
	"github.com/portapps/portapps/pkg/utl"
	"github.com/portapps/portapps/pkg/win"
	"golang.org/x/sys/windows/registry"
)

type config struct {
	Silent bool `yaml:"silent" mapstructure:"silent"`
}

var (
	app *App
	cfg *config
)

func init() {
	var err error

	// Default config
	cfg = &config{
		Silent: false,
	}

	// Init app
	if app, err = NewWithCfg("gnupg-portable", "GnuPG", cfg); err != nil {
		Log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	var err error
	var resp int

	gnupgHome := utl.CreateFolder(utl.PathJoin(app.DataPath, ".gnupg"))

	if !cfg.Silent {
		resp, err = dialog.MsgBox(
			fmt.Sprintf("%s portable", app.Name),
			"Would you like to set GNUPGHOME in your environment ?",
			dialog.MsgBoxBtnYesNo|dialog.MsgBoxIconQuestion)
		if err != nil {
			Log.Fatal().Err(err).Msg("Cannot create dialog box")
		}
	} else {
		resp = dialog.MsgBoxSelectYes
	}

	if resp != dialog.MsgBoxSelectYes {
		Log.Info().Msg("Skipping setting GNUPGHOME...")
		return
	}

	Log.Info().Msgf("Set GNUPGHOME=%s", gnupgHome)
	err = win.SetPermEnv(registry.CURRENT_USER, "GNUPGHOME", gnupgHome)
	if err != nil {
		Log.Fatal().Err(err).Msg("Cannot set GNUPGHOME")
	}
	win.RefreshEnv()
}
