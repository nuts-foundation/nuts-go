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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go/types"
	"os"
	"reflect"
	"strings"
)

const defaultPrefix = "NUTS"
const defaultSeparator = "."
const defaultConfigFile = "nuts.yaml"
const configFileFlag = "configfile"

var defaultIgnoredPrefixes = []string{"root"}

// NutsGlobalConfig has the settings which influence all other settings.
type NutsGlobalConfig struct {
	// The default config file the configuration looks for (Default nuts.yaml)
	DefaultConfigFile string

	// Prefix sets the global config environment variable prefix (Default: NUTS)
	Prefix string

	// Delimiter sets the nested config separator string (Default: '.')
	Delimiter string

	// IgnoredPrefixes is a slice of prefixes which will not be used to prepend config variables, eg: --logging.verbosity will just be --verbosity
	IgnoredPrefixes []string

	v *viper.Viper
}

// NewNutsGlobalConfig creates a NutsGlobalConfig with the following defaults
// * Prefix: NUTS
// * Delimiter: '.'
// * IgnoredPrefixes: ["root","logging"]
func NewNutsGlobalConfig() *NutsGlobalConfig {
	return &NutsGlobalConfig{
		DefaultConfigFile: defaultConfigFile,
		Prefix:            defaultPrefix,
		Delimiter:         defaultSeparator,
		IgnoredPrefixes:   defaultIgnoredPrefixes,
		v:                 viper.New(),
	}
}

// Load sets some initial config in order to be able for commands to load the right parameters and to add the configFile Flag.
// This is mainly spf13/viper related stuff
func (ngc *NutsGlobalConfig) Load() error {
	ngc.v.SetEnvPrefix(ngc.Prefix)
	ngc.v.AutomaticEnv()
	ngc.v.SetEnvKeyReplacer(strings.NewReplacer(ngc.Delimiter, "_"))
	flagSet := pflag.NewFlagSet("config", pflag.ContinueOnError)
	flagSet.String(configFileFlag, ngc.DefaultConfigFile, "Nuts config file")
	pflag.CommandLine.AddFlagSet(flagSet)

	cf := flagSet.Lookup(configFileFlag)

	if err := ngc.v.BindPFlag(cf.Name, cf); err != nil {
		return err
	}

	if err := ngc.v.BindEnv(cf.Name); err != nil {
		return err
	}

	// load flags into viper
	pflag.Parse()

	return ngc.loadConfigFile()
}

// LoadConfigFile load the config from the given config file or from the default config file. If the file does not exist it'll continue with default values.
func (ngc *NutsGlobalConfig) loadConfigFile() error {
	// first load configFile param
	if !ngc.v.IsSet(configFileFlag) {
		return types.Error{Msg: "no configFile is set, run Load before running LoadConfigFile"}
	}
	configFile := ngc.v.GetString(configFileFlag)

	// default path, relative paths and absolute paths should work
	ngc.v.AddConfigPath(".")
	ngc.v.SetConfigFile(configFile)

	// if file can not be found, print to stderr and continue
	err := ngc.v.ReadInConfig()
	if err != nil && err.Error() == fmt.Sprintf("open %s: no such file or directory", configFile) {
		fmt.Fprintf(os.Stderr, "Config file %s not found, using defaults!\n", configFile)
		return nil
	}
	return err
}

// InjectIntoEngine loop over all flags from an engine and injects any value into the given Config struct for the Engine.
// If the Engine does not have a config struct, it does nothing.
// Any config not registered as global flag will be ignored.
// It expects all config var names to be prepended or nested with the Engine ConfigKey,
// this will be ignored if the ConfigKey is "" or if the key is in the set of ignored prefixes.
func (ngc *NutsGlobalConfig) InjectIntoEngine(e *Engine) error {
	var err error

	// todo: trace logging

	// ignore if no target for injection
	if e.Config != nil {
		// ignore if no registered flags
		if e.FlagSet != nil {
			fs := e.FlagSet

			fs.VisitAll(func(f *pflag.Flag) {
				// config name as used by viper
				configName := ngc.configName(e, f)

				// field in struct
				var field *reflect.Value
				field, err = ngc.findField(e, ngc.fieldName(e, f.Name))

				if err != nil {
					err = types.Error{Msg: fmt.Sprintf("Problem injecting [%v] for %s: %s", configName, e.Name, err.Error())}
					return
				}

				// get value
				val := ngc.v.Get(configName)

				if val == nil {
					err = types.Error{Msg: fmt.Sprintf("Nil value for %v, forgot to add flag binding?", configName)}
					return
				}

				// inject value
				field.Set(reflect.ValueOf(val))
			})
		}
	}

	return err
}

func (ngc *NutsGlobalConfig) injectIntoStruct(s interface{}) error {
	var err error

	for _, configName := range ngc.v.AllKeys() {
		// ignore configFile flag
		if configName == configFileFlag {
			continue
		}

		sv := reflect.ValueOf(s)
		var field *reflect.Value
		field, err = ngc.findFieldInStruct(&sv, configName)

		if err != nil {
			return types.Error{Msg: fmt.Sprintf("Problem injecting [%v]: %s", configName, err.Error())}
		}

		// get value
		val := ngc.v.Get(configName)

		if val == nil {
			return types.Error{Msg: fmt.Sprintf("Nil value for %v, forgot to add flag binding?", configName)}
		}

		// inject value
		field.Set(reflect.ValueOf(val))
	}
	return err
}

// RegisterFlags adds the flagSet of an engine to the commandline, flag names are prefixed if needed
func (ngc *NutsGlobalConfig) RegisterFlags(e *Engine) {
	if e.FlagSet != nil {
		fs := e.FlagSet

		fs.VisitAll(func(f *pflag.Flag) {
			// prepend with engine.configKey
			if e.ConfigKey != "" && !ngc.isIgnoredPrefix(e.ConfigKey) {
				f.Name = fmt.Sprintf("%s%s%s", e.ConfigKey, ngc.Delimiter, f.Name)
			}

			// add commandline flag
			pf := pflag.CommandLine.Lookup(f.Name)
			if pf == nil {
				pflag.CommandLine.AddFlag(f)
				pf = f
			}

			// some magic for stuff to get combined
			ngc.v.BindPFlag(f.Name, pf)

			// bind environment variable
			ngc.v.BindEnv(f.Name)
		})
	}
}

func (ngc *NutsGlobalConfig) isIgnoredPrefix(prefix string) bool {
	for _, ip := range ngc.IgnoredPrefixes {
		if ip == prefix {
			return true
		}
	}
	return false
}

// Unmarshal loads config from Env, commandLine and configFile into given struct.
// This call is intended to be used outside of the engine structure of Nuts-go.
// It can be used by the individual repo's, for testing the repo as standalone command.
func (ngc *NutsGlobalConfig) LoadAndUnmarshal(targetCfg interface{}) error {
	if err := ngc.Load(); err != nil {
		return err
	}

	return ngc.injectIntoStruct(targetCfg)
}

// configName returns the fully qualified config name including prefixes and delimiter
func (ngc *NutsGlobalConfig) configName(e *Engine, f *pflag.Flag) string {
	if e.ConfigKey == "" {
		return f.Name
	}
	for _, i := range ngc.IgnoredPrefixes {
		if i == e.ConfigKey {
			return f.Name
		}
	}

	// check if flag name already starts with prefix
	if strings.Index(f.Name, e.ConfigKey) == 0 {
		return f.Name
	}

	// add prefix
	return fmt.Sprintf("%s%s%s", e.ConfigKey, ngc.Delimiter, f.Name)
}

func (ngc *NutsGlobalConfig) fieldName(e *Engine, s string) string {
	if e.ConfigKey != "" && !ngc.isIgnoredPrefix(e.ConfigKey) {
		if strings.Index(s, e.ConfigKey) == 0 {
			return s[len(e.ConfigKey)+1:]
		}
	}

	return s
}

// findField returns the Value of the field to inject value into
// it also checks if the Field can be set
// it uses findFieldRecursive to find deeper nested struct fields
func (ngc *NutsGlobalConfig) findField(e *Engine, fieldName string) (*reflect.Value, error) {
	cfgP := reflect.ValueOf(e.Config)

	return ngc.findFieldInStruct(&cfgP, fieldName)
}

func (ngc *NutsGlobalConfig) findFieldInStruct(cfgP *reflect.Value, configName string) (*reflect.Value, error) {
	if cfgP.Kind() != reflect.Ptr {
		return nil, types.Error{Msg: "Only struct pointers are supported to be a Config target"}
	}

	s := cfgP.Elem()
	if !s.CanSet() {
		return nil, types.Error{Msg: "Given Engine.Config can not be Altered"}
	}

	spl := strings.Split(configName, ngc.Delimiter)

	return ngc.findFieldRecursive(&s, spl)
}

func (ngc *NutsGlobalConfig) findFieldRecursive(s *reflect.Value, names []string) (*reflect.Value, error) {
	head := names[0]
	tail := names[1:]

	t := strings.Title(head)
	field := s.FieldByName(t)
	switch field.Kind() {
	case reflect.Invalid:
		return nil, types.Error{Msg: fmt.Sprintf("inaccessible or invalid field [%v] in %v", t, s.Type())}
	case reflect.Struct:
		if len(tail) == 0 {
			return nil, types.Error{Msg: fmt.Sprintf("incompatible source/target, trying to set value to struct target: %v to %v", strings.Title(head), field.Type())}
		}
		return ngc.findFieldRecursive(&field, tail)
	case reflect.Map:
		return nil, types.Error{Msg: fmt.Sprintf("Map values not supported in %v", field.Type())}
	default:
		if len(tail) > 0 {
			n := fmt.Sprintf("%s.%s", head, strings.Join(tail, "."))
			return nil, types.Error{Msg: fmt.Sprintf("incompatible source/target, deeper nested key than target %s", n)}
		}
	}

	if !field.CanSet() {
		return nil, types.Error{Msg: fmt.Sprintf("Field %v can not be Set", t)}
	}

	return &field, nil
}
