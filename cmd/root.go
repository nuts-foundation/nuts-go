/*
 * Nuts go
 * Copyright (C) 2019 Nuts community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package cmd

import (
	"github.com/labstack/echo/v4"
	consent "github.com/nuts-foundation/nuts-consent-store/engine"
	crypto "github.com/nuts-foundation/nuts-crypto/engine"
	validation "github.com/nuts-foundation/nuts-fhir-validation/engine"
	"github.com/nuts-foundation/nuts-go/pkg"
	registry "github.com/nuts-foundation/nuts-registry/engine"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nuts",
	Short: "The Nuts service executable",
	Run: func(cmd *cobra.Command, args []string) {

		// start interfaces
		echo := echo.New()

		for _, engine := range pkg.EngineCtl.Engines {
			engine.Routes(echo)
		}

		defer shutdownEngines()

		logrus.Fatal(echo.Start("localhost:5678"))
	},
}

func Execute() {
	//flag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	// register static set of engines, needed for other commands
	registerEngines()

	// add all commands from registered engines
	addSubCommands(rootCmd)

	// Load global Nuts config
	cfg := pkg.NutsConfig()

	// todo: combine the following 3 calls into 1 passing an array of engines
	// add commandline options and parse commandline
	addFlagSets(cfg)

	// Load all config and add generic options
	if err := cfg.Load(); err != nil {
		panic(err)
	}

	// logger is initialized
	cfg.PrintConfig()

	// Load config into engines
	injectConfig(cfg)

	// check config on all engines
	configureEngines()

	// start engines
	startEngines()

	// blocking main call
	rootCmd.Execute()
}

func addSubCommands(root *cobra.Command) {
	for _, e := range pkg.EngineCtl.Engines {
		root.AddCommand(e.Cmd)
	}
}

func registerEngines() {
	pkg.RegisterEngine(crypto.NewCryptoEngine())
	pkg.RegisterEngine(validation.NewValidationEngine())
	pkg.RegisterEngine(consent.NewConsentStoreEngine())
	pkg.RegisterEngine(registry.NewRegistryEngine())
}

func injectConfig(cfg *pkg.NutsGlobalConfig) {
	// loop through configs and call viper.Get prepended with engine ConfigKey, inject value into struct
	for _, e := range pkg.EngineCtl.Engines {
		if err := cfg.InjectIntoEngine(e); err != nil {
			logrus.Fatal(err)
		}
	}
}

func configureEngines() {
	for _, e := range pkg.EngineCtl.Engines {
		// only if Engine is dynamically configurable
		if e.Configure != nil {
			if err := e.Configure(); err != nil {
				logrus.Fatal(err)
			}
		}
	}
}

func addFlagSets(cfg *pkg.NutsGlobalConfig) {
	for _, e := range pkg.EngineCtl.Engines {
		cfg.RegisterFlags(e)
	}
}

func startEngines() {
	for _, e := range pkg.EngineCtl.Engines {
		if e.Start != nil {
			if err := e.Start(); err != nil {
				logrus.Fatal(err)
			}
		}
	}
}

func shutdownEngines() {
	for _, e := range pkg.EngineCtl.Engines {
		if e.Shutdown != nil {
			if err := e.Shutdown(); err != nil {
				logrus.Error(err)
			}
		}
	}
}
