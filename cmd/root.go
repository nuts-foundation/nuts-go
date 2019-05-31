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
	goflag "flag"
	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-crypto/pkg/crypto"
	"github.com/nuts-foundation/nuts-fhir-validation/pkg/validation"

	//"github.com/nuts-foundation/nuts-crypto/pkg/crypto"
	//"github.com/nuts-foundation/nuts-fhir-validation/pkg/validation"
	"github.com/nuts-foundation/nuts-go/pkg"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "nuts",
	Short: "The Nuts service executable",
	Run: func(cmd *cobra.Command, args []string) {

		// start engines & monitoring

		// start interfaces
		echo := echo.New()

		for _, engine := range pkg.EngineCtl.Engines {
			engine.Routes(echo)
		}

		echo.Logger.Fatal(echo.Start("localhost:5678"))
	},
}

func Execute() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	registerEngines()
	configureEngines()
	addSubCommands(rootCmd)
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
}

func configureEngines() {
	for _, e := range pkg.EngineCtl.Engines {
		if e.Configure != nil {
			if err := e.Configure(); err != nil {
				panic(err)
			}
		}
	}
}
