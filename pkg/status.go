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
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
	core "github.com/nuts-foundation/nuts-go-core"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
)

//NewStatusEngine creates a new Engine for viewing all engines
func NewStatusEngine() *core.Engine {
	return &core.Engine{
		Name: "Status",
		Cmd: &cobra.Command{
			Use:   "engineStatus",
			Short: "show the registered engines",
			Run: func(cmd *cobra.Command, args []string) {
				names := listAllEngines()
				fmt.Println(strings.Join(names, ","))
			},
		},
		Routes: func(router runtime.EchoRouter) {
			router.GET("/status/engines", ListAllEngines)
		},
	}
}

func init() {
	core.RegisterEngine(NewStatusEngine())
}

// ListAllEngines is the handler function for the /status/engines api call
func ListAllEngines(ctx echo.Context) error {
	names := listAllEngines()

	// generate output
	return ctx.JSON(http.StatusOK, names)
}

func listAllEngines() []string {
	var names []string
	for _, e := range core.EngineCtl.Engines {
		names = append(names, e.Name)
	}
	return names
}
