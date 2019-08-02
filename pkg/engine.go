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
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

	// ConfigKey is the root yaml key in the config file or ENV sub-key for all keys used to configure an engine
	// 	status:
	//	  key:
	// and
	// 	NUTS_STATUS_KEY=
	// and
	//	--status-key
	ConfigKey string

	// Config is the pointer to a config struct. The config will be unmarshalled using the ConfigKey.
	Config interface{}

	// Configure checks if the combination of config parameters is allowed
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
}
