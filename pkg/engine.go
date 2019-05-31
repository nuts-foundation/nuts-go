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
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"net/http"
	"strings"
)

// EngineCtl is the control structure where engines are registered. All registered engines are referenced by the EngineCtl
type EngineControl struct {
	// Engines is the slice of all registered engines
	Engines []*Engine
}

var EngineCtl EngineControl

// Engine contains all the configuration options and callbacks needed by the executable to configure, start, monitor and shutdown the engines
type Engine struct {
	// Name holds the human readable name of the engine
	Name string

	// Cmd is the optional sub-command for the engine. An engine can only add one sub-command (but multiple sub-sub-commands for the sub-command)
	Cmd *cobra.Command

	// Configure loads the given configurations in the engine. Any wrong combination will return an error
	Configure func() error

	// FlasSet contains all engine-local configuration possibilities so they can be displayed through the help command
	FlagSet *pflag.FlagSet

	// Routes passes the Echo router to the specific engine for it to register their routes.
	Routes func(router runtime.EchoRouter)

	// Shutdown the engine
	Shutdown func() error

	// Start the engine, this will spawn any clients, background tasks or active processes.
	Start func() error
}

// RegisterEngine is a helper func to add an engine to the list of engines from a different lib/pkg
func RegisterEngine(engine *Engine) {
	EngineCtl.registerEngine(engine)
}

func (ec *EngineControl) registerEngine(engine *Engine) {
	ec.Engines = append(ec.Engines, engine)
}

func init() {
	EngineCtl = EngineControl{}
	EngineCtl.registerEngine(NewStatusEngine())
}

//NewStatusEngine creates a new Engine for viewing all engines
func NewStatusEngine() *Engine {
	return &Engine{
		Cmd: &cobra.Command{
			Use:   "engineStatus",
			Short: "show the registered engines",
			Run: func(cmd *cobra.Command, args []string) {
				names := listAllEngines()
				fmt.Println(strings.Join(names, ","))
			},
		},
		Name: "Status",
		Routes: func(router runtime.EchoRouter) {
			router.GET("/status/engines", ListAllEngines)
		},
	}
}

// ListAllEngines is the handler function for the /status/engines api call
func ListAllEngines(ctx echo.Context) error {
	names := listAllEngines()

	// generate output
	return ctx.JSON(http.StatusOK, names)
}

func listAllEngines() []string {
	var names []string
	for _, e := range EngineCtl.Engines{
		names = append(names, e.Name)
	}
	return names
}
