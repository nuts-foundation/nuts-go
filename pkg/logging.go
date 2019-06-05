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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)


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
	}
}

func init() {
	EngineCtl.registerEngine(NewLoggerEngine())
}

func printLoggerSetup() {
	log.Infof("Verbosity is set to %s\n", NutsConfig().v.GetString(loggerLevelFlag))
}
