//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"fmt"

	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/log"
	"github.com/portapps/portapps/v3/pkg/utl"
	"github.com/portapps/portapps/v3/pkg/win"
	"golang.org/x/sys/windows/registry"
)

type config struct {
	Silent bool `yaml:"silent" mapstructure:"silent"`
}

var (
	app *portapps.App
	cfg *config
)

func init() {
	var err error

	// Default config
	cfg = &config{
		Silent: false,
	}

	// Init app
	if app, err = portapps.NewWithCfg("gnupg-portable", "GnuPG", cfg); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	var err error
	var resp int

	gnupgHome := utl.CreateFolder(utl.PathJoin(app.DataPath, ".gnupg"))

	if !cfg.Silent {
		resp, err = win.MsgBox(
			fmt.Sprintf("%s portable", app.Name),
			"Would you like to set GNUPGHOME in your environment ?",
			win.MsgBoxBtnYesNo|win.MsgBoxIconQuestion)
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot create dialog box")
		}
	} else {
		resp = win.MsgBoxSelectYes
	}

	if resp != win.MsgBoxSelectYes {
		log.Info().Msg("Skipping setting GNUPGHOME...")
		return
	}

	log.Info().Msgf("Set GNUPGHOME=%s", gnupgHome)
	err = win.SetPermEnv(registry.CURRENT_USER, "GNUPGHOME", gnupgHome)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot set GNUPGHOME")
	}
	win.RefreshEnv()

	defer app.Close()
}
