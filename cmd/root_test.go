package cmd

import (
	core "github.com/nuts-foundation/nuts-go-core"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_rootCmd(t *testing.T) {
	t.Run("start in CLI mode", func(t *testing.T) {
		var routesCalled = false
		core.RegisterEngine(&core.Engine{
			Routes: func(router core.EchoRouter) {
				routesCalled = true
			},
		})
		os.Setenv("NUTS_MODE", core.GlobalCLIMode)
		defer os.Unsetenv("NUTS_MODE")
		assert.NoError(t, rootCmd.Execute())
		assert.False(t, routesCalled, "engine.Routes was called")
	})
}