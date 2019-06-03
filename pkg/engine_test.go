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
	"github.com/golang/mock/gomock"
	"github.com/nuts-foundation/nuts-go/mock"
	"net/http"
	"testing"
)

type dummyEngine struct {

}

func TestRegisterEngine(t *testing.T) {
	t.Run("adds an engine to the list", func(t *testing.T) {
		ctl := EngineControl{
			Engines: []*Engine{},
		}
		ctl.registerEngine(&Engine{})

		if len(ctl.Engines) != 1 {
			t.Errorf("Expected 1 registered engine, Got %d", len(ctl.Engines))
		}
	})

	t.Run("has been called by init to register StatusEngine", func(t *testing.T) {

		if len(EngineCtl.Engines) != 2 {
			t.Errorf("Expected 2 registered engine, Got %d", len(EngineCtl.Engines))
		}
	})
}

func TestNewStatusEngine_Routes(t *testing.T) {
	t.Run("Registers a single route for listing all engines", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockEchoRouter(ctrl)

		echo.EXPECT().GET("/status/engines", gomock.Any())

		NewStatusEngine().Routes(echo)
	})
}

func TestNewStatusEngine_Cmd(t *testing.T) {
	t.Run("Cmd returns a cobra command", func(t *testing.T) {
		e := NewStatusEngine().Cmd
		if e.Name() != "engineStatus" {
			t.Errorf("Expected a command with name engineStatus, Got %s", e.Name())
		}
	})
}

func TestListAllEngines(t *testing.T) {
	t.Run("ListAllEngines renders json output of list of engines", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		echo.EXPECT().JSON(http.StatusOK, []string{"Logging", "Status"})

		ListAllEngines(echo)
	})
}