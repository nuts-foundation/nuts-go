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
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"math"
	"os"
	"reflect"
	"strings"
	"sync"
)

const defaultPrefix = "NUTS"
const defaultSeparator = "."
const defaultConfigFile = "nuts.yaml"
const configFileFlag = "configfile"
const loggerLevelFlag = "verbosity"
const defaultLogLevel = "info"

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

var configOnce sync.Once
var configInstance *NutsGlobalConfig

// NutsGlobalConfig returns a singleton global config
func NutsConfig() *NutsGlobalConfig {
	configOnce.Do(func() {
		configInstance = NewNutsGlobalConfig()
	})
	return configInstance
}

// Load sets some initial config in order to be able for commands to load the right parameters and to add the configFile Flag.
// This is mainly spf13/viper related stuff
func (ngc *NutsGlobalConfig) Load(cmd *cobra.Command) error {
	ngc.v.SetEnvPrefix(ngc.Prefix)
	ngc.v.AutomaticEnv()
	ngc.v.SetEnvKeyReplacer(strings.NewReplacer(ngc.Delimiter, "_"))
	flagSet := pflag.NewFlagSet("config", pflag.ContinueOnError)
	flagSet.String(configFileFlag, ngc.DefaultConfigFile, "Nuts config file")
	flagSet.String(loggerLevelFlag, defaultLogLevel, "Log level")
	cmd.PersistentFlags().AddFlagSet(flagSet)

	// Bind config flag
	// Bind log level flag
	ngc.bindFlag(flagSet, configFileFlag)
	ngc.bindFlag(flagSet, loggerLevelFlag)

	// load flags into viper
	pfs := cmd.PersistentFlags()
	pfs.ParseErrorsWhitelist.UnknownFlags = true
	if err := pfs.Parse(os.Args[1:]); err != nil {
		if err != pflag.ErrHelp {
			return err
		}
	}

	// load configFile into viper
	if err := ngc.loadConfigFile(); err != nil {
		return err
	}

	// initialize logger, verbosity flag needs to be available
	level, err := log.ParseLevel(ngc.v.GetString(loggerLevelFlag))
	if err != nil {
		return err
	}
	log.SetLevel(level)

	return nil
}

func (ngc *NutsGlobalConfig) bindFlag(fs *pflag.FlagSet, name string) error {
	s := fs.Lookup(name)
	if err := ngc.v.BindPFlag(s.Name, s); err != nil {
		return err
	}
	if err := ngc.v.BindEnv(s.Name); err != nil {
		return err
	}
	return nil
}

// PrintConfig outputs the current config to the logger on info level
func (ngc *NutsGlobalConfig) PrintConfig(logger log.FieldLogger) {
	title := "Config"
	var longestKey = 10
	var longestValue int
	for _, e := range EngineCtl.Engines {
		if e.FlagSet != nil {
			e.FlagSet.VisitAll(func(flag *pflag.Flag) {
				s := fmt.Sprintf("%v", ngc.v.Get(strings.ToLower(flag.Name)))
				if len(s) > longestValue {
					longestValue = len(s)
				}
				if len(flag.Name) > longestKey {
					longestKey = len(flag.Name)
				}
			})
		}
	}

	totalLength := 7 + longestKey + longestValue
	stars := strings.Repeat("*", totalLength)
	sideStarsLeft := int(math.Floor((float64(totalLength)-float64(len(title)))/2.0)) - 1
	sideStarsRight := int(math.Ceil((float64(totalLength)-float64(len(title)))/2.0)) - 1

	logger.Infoln(stars)
	logger.Infof("%s %s %s", strings.Repeat("*", sideStarsLeft), title, strings.Repeat("*", sideStarsRight))

	f := fmt.Sprintf("%%-%ds%%v", 7+longestKey)

	logger.Infof(f, configFileFlag, ngc.v.Get(configFileFlag))
	logger.Infof(f, loggerLevelFlag, ngc.v.Get(loggerLevelFlag))
	for _, e := range EngineCtl.Engines {
		if e.FlagSet != nil {
			e.FlagSet.VisitAll(func(flag *pflag.Flag) {
				logger.Infof(f, flag.Name, ngc.v.Get(strings.ToLower(flag.Name)))
			})
		}
	}

	logger.Infoln(stars)
}

// LoadConfigFile load the config from the given config file or from the default config file. If the file does not exist it'll continue with default values.
func (ngc *NutsGlobalConfig) loadConfigFile() error {
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

	// ignore if no target for injection
	if e.Config != nil {
		// ignore if no registered flags
		if e.FlagSet != nil {
			fs := e.FlagSet
			log.Tracef("Injecting values for engine %s\n", e.Name)

			fs.VisitAll(func(f *pflag.Flag) {
				// config name as used by viper
				configName := ngc.configName(e, f)

				// field in struct
				var field *reflect.Value
				field, err = ngc.findField(e, ngc.fieldName(e, f.Name))

				if err != nil {
					err = errors.New(fmt.Sprintf("Problem injecting [%v] for %s: %s", configName, e.Name, err.Error()))
					return
				}

				// get value
				val := ngc.v.Get(configName)

				if val == nil {
					err = errors.New(fmt.Sprintf("Nil value for %v, forgot to add flag binding?", configName))
					return
				}

				// inject value
				field.Set(reflect.ValueOf(val))
				log.Tracef("[%s] %s=%v\n", e.Name, f.Name, val)
			})
		}
	}

	return err
}

func (ngc *NutsGlobalConfig) injectIntoStruct(s interface{}) error {
	var err error

	for _, configName := range ngc.v.AllKeys() {
		// ignore configFile flag
		if configName == configFileFlag || configName == loggerLevelFlag {
			continue
		}

		sv := reflect.ValueOf(s)
		var field *reflect.Value
		field, err = ngc.findFieldInStruct(&sv, configName)

		if err != nil {
			return errors.New(fmt.Sprintf("Problem injecting [%v]: %s", configName, err.Error()))
		}

		// inject value
		field.Set(reflect.ValueOf(ngc.v.Get(configName)))
	}
	return err
}

// RegisterFlags adds the flagSet of an engine to the commandline, flag names are prefixed if needed
// The passed command must be the root command not the engine.Cmd (unless they are the same)
func (ngc *NutsGlobalConfig) RegisterFlags(cmd *cobra.Command, e *Engine) {
	if e.FlagSet != nil {
		fs := e.FlagSet

		fs.VisitAll(func(f *pflag.Flag) {
			// prepend with engine.configKey
			if e.ConfigKey != "" && !ngc.isIgnoredPrefix(e.ConfigKey) {
				f.Name = fmt.Sprintf("%s%s%s", e.ConfigKey, ngc.Delimiter, f.Name)
			}

			// add commandline flag
			pf := cmd.PersistentFlags().Lookup(f.Name)
			if pf == nil {
				cmd.PersistentFlags().AddFlag(f)
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
func (ngc *NutsGlobalConfig) LoadAndUnmarshal(cmd *cobra.Command, targetCfg interface{}) error {
	if err := ngc.Load(cmd); err != nil {
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
		return nil, errors.New("Only struct pointers are supported to be a Config target")
	}

	s := cfgP.Elem()
	if !s.CanSet() {
		return nil, errors.New("Given Engine.Config can not be Altered")
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
		return nil, errors.New(fmt.Sprintf("inaccessible or invalid field [%v] in %v", t, s.Type()))
	case reflect.Struct:
		if len(tail) == 0 {
			return nil, errors.New(fmt.Sprintf("incompatible source/target, trying to set value to struct target: %v to %v", strings.Title(head), field.Type()))
		}
		return ngc.findFieldRecursive(&field, tail)
	case reflect.Map:
		return nil, errors.New(fmt.Sprintf("Map values not supported in %v", field.Type()))
	default:
		if len(tail) > 0 {
			n := fmt.Sprintf("%s.%s", head, strings.Join(tail, "."))
			return nil, errors.New(fmt.Sprintf("incompatible source/target, deeper nested key than target %s", n))
		}
	}

	if !field.CanSet() {
		return nil, errors.New(fmt.Sprintf("Field %v can not be Set", t))
	}

	return &field, nil
}
