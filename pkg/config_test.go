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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"testing"
)

func TestNewNutsGlobalConfig(t *testing.T) {
	t.Run("returns a NutsGlobalConfig with defaults", func(t *testing.T) {
		c := NewNutsGlobalConfig()

		if c.DefaultConfigFile != defaultConfigFile {
			t.Errorf("Expected DefaultConfigFile to be [%s], got [%s]", defaultConfigFile, c.DefaultConfigFile)
		}

		if c.Prefix != defaultPrefix {
			t.Errorf("Expected Prefix to be [%s], got [%s]", defaultPrefix, c.Prefix)
		}

		if c.Delimiter != defaultSeparator {
			t.Errorf("Expected Prefix to be [%s], got [%s]", defaultSeparator, c.Delimiter)
		}

		if !reflect.DeepEqual(c.IgnoredPrefixes, defaultIgnoredPrefixes) {
			t.Errorf("Expected Prefix to be [%s], got [%s]", defaultIgnoredPrefixes, c.IgnoredPrefixes)
		}
	})
}

func TestNutsGlobalConfig_Configure(t *testing.T) {
	if err := NewNutsGlobalConfig().Configure(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	t.Run("Sets global Env prefix", func(t *testing.T) {
		os.Setenv("NUTS_KEY", "value")
		if value := viper.Get("key"); value != "value" {
			t.Errorf("Expected key to have [value], got [%v]", value)
		}
	})

	t.Run("Sets correct key replacer", func(t *testing.T) {
		os.Setenv("NUTS_SUB_KEY", "value")
		if value := viper.Get("sub.key"); value != "value" {
			t.Errorf("Expected sub.key to have [value], got [%v]", value)
		}
	})

	t.Run("Adds configFile flag", func(t *testing.T) {
		if value := viper.Get(configFileFlag); value != defaultConfigFile {
			t.Errorf("Expected configFile to be [%s], got [%v]", defaultConfigFile, value)
		}
	})
}

func TestNutsGlobalConfig_LoadConfigFile(t *testing.T) {
	t.Run("Does not return error on missing file", func(t *testing.T) {
		cfg := NutsGlobalConfig{
			DefaultConfigFile: "non_existing.yaml",
		}
		cfg.Configure()

		if err := cfg.LoadConfigFile(); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Returns error on incorrect file", func(t *testing.T) {
		cfg := NutsGlobalConfig{
			DefaultConfigFile: "../test/config/corrupt.yaml",
		}
		cfg.Configure()

		err := cfg.LoadConfigFile()
		if err == nil {
			t.Errorf("Expected error, got nothing")
		}

		expected := "While parsing config: yaml: line 1: did not find expected ',' or '}'"
		if err.Error() != expected {
			t.Errorf("Expected error: [%s], got [%v]", expected, err.Error())
		}
	})

	t.Run("Loads settings into viper", func(t *testing.T) {
		cfg := NutsGlobalConfig{
			DefaultConfigFile: "../test/config/dummy.yaml",
		}
		cfg.Configure()

		err := cfg.LoadConfigFile()
		if err != nil {
			t.Errorf("Expected no error, got [%v]", err.Error())
		}

		val := viper.Get("key")
		if val != "value" {
			t.Errorf("Expected value to equals [value], got [%v]", val)
		}
	})
}

func TestNutsGlobalConfig_LoadAndUnmarshal(t *testing.T) {
	cfg := NewNutsGlobalConfig()
	cfg.Configure()

	t.Run("Adds configFile flag to Cmd", func(t *testing.T) {
		err := cfg.LoadAndUnmarshal(&struct{}{})
		if err != nil {
			t.Errorf("Expected no error, got [%v]", err.Error())
		}

		if !viper.IsSet(configFileFlag) {
			t.Errorf("Expected %s to be set", configFileFlag)
		}
	})

	t.Run("Injects custom config into struct", func(t *testing.T) {
		s := struct{
			Key string
		}{
			"",
		}
		viper.Set("key", "value")
		err := cfg.LoadAndUnmarshal(&s)
		if err != nil {
			t.Errorf("Expected no error, got [%v]", err.Error())
		}

		if s.Key != "value" {
			t.Errorf("Expected value for key in struct to equals [value], got %s", s.Key)
		}
	})
}

func TestNutsGlobalConfig_InjectIntoEngine(t *testing.T) {
	cfg := NewNutsGlobalConfig()
	cfg.Configure()

	t.Run("param is injected into engine without ConfigKey", func(t *testing.T) {
		c := struct {
			Key string
		}{}

		e := &Engine{
			Config: &c,
			FlagSet: pflag.NewFlagSet("dummy", pflag.ContinueOnError),
		}
		e.FlagSet.String("key", "", "")

		viper.Set("key", "value")

		if err := cfg.InjectIntoEngine(e); err != nil {
			t.Errorf("Expected no error, got [%v]", err.Error())
		}

		if c.Key != "value" {
			t.Errorf("Expected value to be injected into struct")
		}
	})

	t.Run("param is injected into engine with ConfigKey", func(t *testing.T) {
		c := struct {
			Key string
		}{}

		e := &Engine{
			Config: &c,
			ConfigKey: "pre",
			FlagSet: pflag.NewFlagSet("dummy", pflag.ContinueOnError),
		}
		e.FlagSet.String("key", "", "")
		cfg.RegisterFlags(e)

		viper.Set("pre.key", "value")

		if err := cfg.InjectIntoEngine(e); err != nil {
			t.Errorf("Expected no error, got [%v]", err.Error())
		}

		if c.Key != "value" {
			t.Errorf("Expected value to be injected into struct")
		}
	})

	t.Run("nested param is injected into engine without ConfigKey", func(t *testing.T) {
		c := struct {
			Nested struct{Key string}
		}{}

		e := &Engine{
			Config: &c,
			FlagSet: pflag.NewFlagSet("dummy", pflag.ContinueOnError),
		}
		e.FlagSet.String("nested.key", "", "")

		viper.Set("nested.key", "value")

		if err := cfg.InjectIntoEngine(e); err != nil {
			t.Errorf("Expected no error, got [%v]", err.Error())
		}

		if c.Nested.Key != "value" {
			t.Errorf("Expected value to be injected into struct")
		}
	})

	t.Run("nested param is injected into engine with ConfigKey", func(t *testing.T) {
		c := struct {
			Nested struct{Key string}
		}{}

		e := &Engine{
			Config: &c,
			ConfigKey: "pre",
			FlagSet: pflag.NewFlagSet("dummy", pflag.ContinueOnError),
		}
		e.FlagSet.String("nested.key", "", "")
		cfg.RegisterFlags(e)

		viper.Set("pre.nested.key", "value")

		if err := cfg.InjectIntoEngine(e); err != nil {
			t.Errorf("Expected no error, got [%v]", err.Error())
		}

		if c.Nested.Key != "value" {
			t.Errorf("Expected value to be injected into struct")
		}
	})

	t.Run("returns error for inaccessible key in struct", func(t *testing.T) {
		c := struct {
			key string
		}{}

		e := &Engine{
			Name: "test",
			Config: &c,
			FlagSet: pflag.NewFlagSet("dummy", pflag.ContinueOnError),
		}
		e.FlagSet.String("key", "", "")

		err := cfg.InjectIntoEngine(e)
		if err == nil {
			t.Errorf("Expected error, got nothing")
		}

		expected := "-: Problem injecting [key] for test: -: inaccessible or invalid field [Key] in struct { key string }"
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err.Error())
		}
	})

	t.Run("returns error on missing default", func(t *testing.T) {
		c := struct {
			Nested struct{Key string}
		}{}

		e := &Engine{
			Config: &c,
			FlagSet: pflag.NewFlagSet("dummy", pflag.ContinueOnError),
		}
		e.FlagSet.String("nested.key", "", "")

		viper.Set("nested", "value")

		err := cfg.InjectIntoEngine(e)
		if err == nil {
			t.Errorf("Expected error, got nothing")
		}

		expected := "-: Nil value for nested.key, forgot to add flag binding?"
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err.Error())
		}
	})

	t.Run("returns error on wrong nesting", func(t *testing.T) {
		c := struct {
			Nested string
		}{}

		e := &Engine{
			Config: &c,
			FlagSet: pflag.NewFlagSet("dummy", pflag.ContinueOnError),
		}
		e.FlagSet.String("nested.key", "", "")

		viper.Set("nested.key", "value")

		err := cfg.InjectIntoEngine(e)
		if err == nil {
			t.Errorf("Expected error, got nothing")
		}

		expected := "-: Problem injecting [nested.key] for : -: incompatible source/target, deeper nested key than target nested.key"
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err.Error())
		}
	})
}
