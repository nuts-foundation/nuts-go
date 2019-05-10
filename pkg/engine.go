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
	"github.com/labstack/echo"
	"github.com/spf13/pflag"
	"net/http"
	"reflect"
)

// EngineCtl is the control structure where engines are registered. All registered engines are referenced by the EngineCtl
type EngineControl struct {
	// Engines is the slice of all registered engines
	Engines []Engine
}

var EngineCtl EngineControl

// EngineAPIRoute is the structure that holds which api call is routed to which go func.
type EngineAPIRoute struct {
	Path    string
	Handler echo.HandlerFunc
}

// Engine contains all the functions needed by the executable to configure, start, monitor and shutdown the engines
type Engine interface {
	// FlasSet returns all configuration possibilities so they can be displayed through the help command
	FlagSet() *pflag.FlagSet

	// Configure loads the given configurations in the engine. Any wrong combination will return an error
	Configure() error

	// ServerHandlerFunctions gives a list of path, handler functions combinations which should be registered to the echo webserver
	Routes() []EngineAPIRoute

	// Start the engine, this will spawn any clients, background tasks or active processes.
	Start() error

	// Shutdown the engine
	Shutdown() error
}

// RegisterEngine is a helper func to add an engine to the list of engines from a different lib/pkg
func RegisterEngine(engine Engine) {
	EngineCtl.registerEngine(engine)
}

func (ec *EngineControl) registerEngine(engine Engine) {
	ec.Engines = append(ec.Engines, engine)
}

// StatusEngine is an engine that comes with the executable and lists installed engines on an endpoint.
type StatusEngine struct {
}

func init() {
	EngineCtl = EngineControl{}
	EngineCtl.registerEngine(&StatusEngine{})
}

// FlagSet returns an empty FlagSet
func (*StatusEngine) FlagSet() *pflag.FlagSet {
	return &pflag.FlagSet{}
}

// Configure does not do anything
func (*StatusEngine) Configure() error {
	return nil
}

// Routes returns a single endpoint listing all available/active engines on /status/engines
func (se *StatusEngine) Routes() []EngineAPIRoute {
	return []EngineAPIRoute{
		{
			Path:    "/status/engines",
			Handler: se.ListAllEngines,
		},
	}
}

// Start does not do anything
func (*StatusEngine) Start() error {
	return nil
}

// Shutdown does not do anything
func (*StatusEngine) Shutdown() error {
	return nil
}

// ListAllEngines is the handler function for the /status/engines api call
func (se *StatusEngine) ListAllEngines(ctx echo.Context) error {
	var names []string
	for _, e := range EngineCtl.Engines {
		names = append(names, reflect.TypeOf(e).String())
	}

	// generate output
	return ctx.JSON(http.StatusOK, names)
}
