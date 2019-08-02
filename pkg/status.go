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
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
)

//NewStatusEngine creates a new Engine for viewing all engines
func NewStatusEngine() *Engine {
	return &Engine{
		Name: "Status",
		Cmd: &cobra.Command{
			Use:   "engineStatus",
			Short: "show the registered engines",
			Run: func(cmd *cobra.Command, args []string) {
				names := listAllEngines()
				fmt.Println(strings.Join(names, ","))
			},
		},
		Routes: func(router EchoRouter) {
			router.GET("/status/engines", ListAllEngines)
		},
	}
}

func init() {
	EngineCtl.registerEngine(NewStatusEngine())
}

// ListAllEngines is the handler function for the /status/engines api call
func ListAllEngines(ctx echo.Context) error {
	names := listAllEngines()

	// generate output
	return ctx.JSON(http.StatusOK, names)
}

func listAllEngines() []string {
	var names []string
	for _, e := range EngineCtl.Engines {
		names = append(names, e.Name)
	}
	return names
}
