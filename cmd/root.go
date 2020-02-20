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
	"github.com/labstack/echo/v4/middleware"
	bridge "github.com/nuts-foundation/consent-bridge-go-client/engine"
	auth "github.com/nuts-foundation/nuts-auth/engine"
	logic "github.com/nuts-foundation/nuts-consent-logic/engine"
	consent "github.com/nuts-foundation/nuts-consent-store/engine"
	crypto "github.com/nuts-foundation/nuts-crypto/engine"
	octopus "github.com/nuts-foundation/nuts-event-octopus/engine"
	validation "github.com/nuts-foundation/nuts-fhir-validation/engine"
	core "github.com/nuts-foundation/nuts-go-core"
	registry "github.com/nuts-foundation/nuts-registry/engine"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nuts",
	Short: "The Nuts service executable",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := core.NutsConfig()
		if cfg.Mode() != core.GlobalServerMode {
			logrus.Error("Please specify a sub command when running in CLI mode")
			_ = cmd.Help()
			return
		}
		// start interfaces
		echo := echo.New()
		echo.HideBanner = true
		echo.Use(middleware.Logger())
		echo.Use(core.DecodeURIPath)

		for _, engine := range core.EngineCtl.Engines {
			if engine.Routes != nil {
				engine.Routes(echo)
			}
		}

		defer shutdownEngines()
		logrus.Fatal(echo.Start(cfg.ServerAddress()))
	},
}

func Execute() {
	//flag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	// register static set of engines, needed for other commands
	registerEngines()

	// add all commands from registered engines
	addSubCommands(rootCmd)

	// Load global Nuts config
	cfg := core.NutsConfig()

	// todo: combine the following 3 calls into 1 passing an array of engines
	// add commandline options and parse commandline
	addFlagSets(rootCmd, cfg)

	// Load all config and add generic options
	if err := cfg.Load(rootCmd); err != nil {
		panic(err)
	}

	// Load config into engines
	injectConfig(cfg)

	cfg.PrintConfig(logrus.StandardLogger())

	// check config on all engines
	configureEngines()

	// start engines
	startEngines()

	// blocking main call
	rootCmd.Execute()
}

func addSubCommands(root *cobra.Command) {
	for _, e := range core.EngineCtl.Engines {
		if e.Cmd != nil {
			root.AddCommand(e.Cmd)
		}
	}
}

func registerEngines() {
	core.RegisterEngine(core.NewStatusEngine())
	core.RegisterEngine(core.NewLoggerEngine())
	core.RegisterEngine(crypto.NewCryptoEngine())
	core.RegisterEngine(registry.NewRegistryEngine())
	core.RegisterEngine(octopus.NewEventOctopusEngine())

	core.RegisterEngine(logic.NewConsentLogicEngine())
	core.RegisterEngine(consent.NewConsentStoreEngine())
	core.RegisterEngine(validation.NewValidationEngine())
	core.RegisterEngine(auth.NewAuthEngine())
	core.RegisterEngine(bridge.NewConsentBridgeClientEngine())
}

func injectConfig(cfg *core.NutsGlobalConfig) {
	// loop through configs and call viper.Get prepended with engine ConfigKey, inject value into struct
	for _, e := range core.EngineCtl.Engines {
		if err := cfg.InjectIntoEngine(e); err != nil {
			logrus.Fatal(err)
		}
	}
}

func configureEngines() {
	for _, e := range core.EngineCtl.Engines {
		// only if Engine is dynamically configurable
		if e.Configure != nil {
			if err := e.Configure(); err != nil {
				logrus.Fatal(err)
			}
		}
	}
}

func addFlagSets(cmd *cobra.Command, cfg *core.NutsGlobalConfig) {
	for _, e := range core.EngineCtl.Engines {
		cfg.RegisterFlags(cmd, e)
	}
}

func startEngines() {
	for _, e := range core.EngineCtl.Engines {
		if e.Start != nil {
			if err := e.Start(); err != nil {
				logrus.Fatal(err)
			}
		}
	}
}

func shutdownEngines() {
	for _, e := range core.EngineCtl.Engines {
		if e.Shutdown != nil {
			if err := e.Shutdown(); err != nil {
				logrus.Error(err)
			}
		}
	}
}
