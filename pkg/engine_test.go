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
	"github.com/golang/mock/gomock"
	mock "github.com/nuts-foundation/nuts-go/mock"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"net/http"
	"testing"
)

type dummyEngine struct {

}

func TestRegisterEngine(t *testing.T) {
	t.Run("adds an engine to the list", func(t *testing.T) {
		ctl := EngineControl{
			Engines: []Engine{},
		}
		ctl.registerEngine(&dummyEngine{})

		if len(ctl.Engines) != 1 {
			t.Errorf("Expected 1 registered engine, Got %d", len(ctl.Engines))
		}
	})

	t.Run("has been called by init to register StatusEngine", func(t *testing.T) {

		if len(EngineCtl.Engines) != 1 {
			t.Errorf("Expected 1 registered engine, Got %d", len(EngineCtl.Engines))
		}
	})
}

func TestStatusEngine_Routes(t *testing.T) {
	t.Run("Registers a single route for listing all engines", func(t *testing.T) {
		se := StatusEngine{}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockEchoRouter(ctrl)

		echo.EXPECT().GET("/status/engines", gomock.Any())

		se.Routes(echo)
	})
}

func TestStatusEngine_Cmd(t *testing.T) {
	t.Run("Cmd returns a cobra command", func(t *testing.T) {
		se := StatusEngine{}

		e := se.Cmd()
		if e.Name() != "engineStatus" {
			t.Errorf("Expected a command with name engineStatus, Got %s", e.Name())
		}
	})
}

func TestStatusEngine_Configure(t *testing.T) {
	t.Run("Configure returns nil", func(t *testing.T) {
		se := StatusEngine{}

		e := se.Configure()
		if e != nil {
			t.Errorf("Expected no error, Got %s", e.Error())
		}
	})
}

func TestStatusEngine_Start(t *testing.T) {
	t.Run("Start returns nil", func(t *testing.T) {
		se := StatusEngine{}

		e := se.Start()
		if e != nil {
			t.Errorf("Expected no error, Got %s", e.Error())
		}
	})
}

func TestStatusEngine_Shutdown(t *testing.T) {
	t.Run("Shutdown returns nil", func(t *testing.T) {
		se := StatusEngine{}

		e := se.Shutdown()
		if e != nil {
			t.Errorf("Expected no error, Got %s", e.Error())
		}
	})
}

func TestStatusEngine_FlagSet(t *testing.T) {
	t.Run("FlagSet returns empty set", func(t *testing.T) {
		se := StatusEngine{}

		e := se.FlagSet()
		if e.HasAvailableFlags() {
			t.Errorf("Expected flagset to be empty")
		}
	})
}

func TestStatusEngine_ListAllEngines(t *testing.T) {
	t.Run("ListAllEngines renders json output of list of engines", func(t *testing.T) {
		se := StatusEngine{}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		echo.EXPECT().JSON(http.StatusOK, []string{"*pkg.StatusEngine"})

		se.ListAllEngines(echo)
	})
}

func (*dummyEngine) FlagSet() *pflag.FlagSet {
	return &pflag.FlagSet{}
}

func (*dummyEngine) Cmd() *cobra.Command {
	return nil
}

func (*dummyEngine) Configure() error {
	return nil
}

func (*dummyEngine) Routes(echoRouter runtime.EchoRouter) {

}

func (*dummyEngine) Start() error {
	return nil
}

func (*dummyEngine) Shutdown() error {
	return nil
}