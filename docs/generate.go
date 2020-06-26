package main

import (
	core "github.com/nuts-foundation/nuts-go-core"
	docs "github.com/nuts-foundation/nuts-go-core/docs"
	"github.com/nuts-foundation/nuts-go/cmd"
	"github.com/spf13/pflag"
	"strings"
)

func main() {
	flags := make(map[string]*pflag.FlagSet)
	command := cmd.CreateCommand()
	core.NutsConfig().Load(command)
	globalFlags := command.PersistentFlags()
	flags[""] = globalFlags
	// Make sure engines are registered
	for _, engine := range core.EngineCtl.Engines {
		if engine.ConfigKey == "" {
			// go-core engine contains global flags and has no config key
			continue
		}
		flagsForEngine := extractFlagsForEngine(engine.ConfigKey, globalFlags)
		if flagsForEngine.HasAvailableFlags() {
			flags[engine.Name] = flagsForEngine
		}
	}
	docs.GeneratePartitionedConfigOptionsDocs("docs/pages/configuration/options.rst", flags)
}

func extractFlagsForEngine(configKey string, flagSet *pflag.FlagSet) *pflag.FlagSet {
	result := pflag.FlagSet{}
	flagSet.VisitAll(func(current *pflag.Flag) {
		if strings.HasPrefix(current.Name, configKey + ".") {
			// This flag belongs to this engine, so copy it and hide it in the input flag set
			flagCopy := *current
			current.Hidden = true
			result.AddFlag(&flagCopy)
		}
	})
	return &result
}
