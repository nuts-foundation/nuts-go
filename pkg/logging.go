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

package pkg

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type loggerConfig struct {
	Verbosity string
}

var logConfig = loggerConfig{}

//NewLoggerEngine creates a new Engine for logging
func NewLoggerEngine() *Engine {
	return &Engine{
		Name: "Logging",
		Cmd: &cobra.Command{
			Use:   "logStatus",
			Short: "show the current logging setup",
			Run: func(cmd *cobra.Command, args []string) {
				printLoggerSetup()
			},
		},
		Config: &logConfig,
		FlagSet: flagSet(),
	}
}

func init() {
	EngineCtl.registerEngine(NewLoggerEngine())
}

func flagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("logging", pflag.ContinueOnError)

	flags.String("verbosity", "info", "logger verbosity: trace, debug, info, warn, error, fatal, panic")

	return flags
}

func printLoggerSetup() {
	fmt.Printf("Verbosity is set to %s\n", logConfig.Verbosity)
}
