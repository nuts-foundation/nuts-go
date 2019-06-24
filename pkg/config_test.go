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
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
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

func TestNutsConfig(t *testing.T) {
	t.Run("returns same instance every time", func(t *testing.T) {
		if NutsConfig() != NutsConfig() {
			t.Error("Expected instance to be the same")
		}
	})
}

func TestNutsGlobalConfig_Load(t *testing.T) {
	cfg := NewNutsGlobalConfig()

	if err := cfg.Load(&cobra.Command{}); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	t.Run("Sets global Env prefix", func(t *testing.T) {
		os.Setenv("NUTS_KEY", "value")
		if value := cfg.v.Get("key"); value != "value" {
			t.Errorf("Expected key to have [value], got [%v]", value)
		}
	})

	t.Run("Sets correct key replacer", func(t *testing.T) {
		os.Setenv("NUTS_SUB_KEY", "value")
		if value := cfg.v.Get("sub.key"); value != "value" {
			t.Errorf("Expected sub.key to have [value], got [%v]", value)
		}
	})

	t.Run("Adds configFile flag", func(t *testing.T) {
		if value := cfg.v.Get(configFileFlag); value != defaultConfigFile {
			t.Errorf("Expected configFile to be [%s], got [%v]", defaultConfigFile, value)
		}
	})
}

func TestNutsGlobalConfig_Load2(t *testing.T) {
	defer func() {
		os.Args = []string{"command"}
	}()

	t.Run("Ignores unknown flags when parsing", func(t *testing.T) {
		os.Args = []string{"executable", "command", "--unknown", "value"}
		cfg := NewNutsGlobalConfig()
		if err := cfg.Load(&cobra.Command{}); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Incorrect arguments returns error", func(t *testing.T) {
		os.Args = []string{"command", "---"}
		cfg := NewNutsGlobalConfig()

		err := cfg.Load(&cobra.Command{})

		if err == nil {
			t.Error("Expected error, got nothing")
			return
		}
	})

	t.Run("Ignores --help as incorrect argument", func(t *testing.T) {
		os.Args = []string{"command", "--help"}
		cfg := NewNutsGlobalConfig()

		if err := cfg.Load(&cobra.Command{}); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})


	t.Run("Returns error for incorrect verbosity", func(t *testing.T) {
		os.Args = []string{"command", "--verbosity", "hell"}
		cfg := NewNutsGlobalConfig()

		err := cfg.Load(&cobra.Command{})

		if err == nil {
			t.Error("Expected error, got nothing")
			return
		}

		expected := "not a valid logrus Level: \"hell\""
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err)
		}
	})
}

func TestNutsGlobalConfig_PrintConfig(t *testing.T) {
	cfg := NewNutsGlobalConfig()
	cfg.v.Set("key", "value")
	logger := logrus.New()
	buf := new(bytes.Buffer)
	logger.Out = buf
	cfg.PrintConfig(logger)
	bs := buf.String()

	t.Run("output contains key", func(t *testing.T) {
		if strings.Index(bs, "key") == -1 {
			t.Error("Expected key to be in output")
		}
	})

	t.Run("output contains some stars", func(t *testing.T) {
		if strings.Index(bs, "***************") == -1 {
			t.Error("Expected stars to be in output")
		}
	})

	t.Run("output contains header", func(t *testing.T) {
		if strings.Index(bs, "*** Config ****") == -1 {
			t.Error("Expected header to be in output")
		}
	})
}

func TestNutsGlobalConfig_RegisterFlags(t *testing.T) {
	t.Run("adds prefix to flag", func(t *testing.T) {
		e := &Engine{
			Cmd:       &cobra.Command{},
			ConfigKey: "pre",
			FlagSet:   pflag.NewFlagSet("dummy", pflag.ContinueOnError),
		}
		e.FlagSet.String("key", "", "")

		cfg := NewNutsGlobalConfig()
		cfg.RegisterFlags(e.Cmd, e)

		var found bool
		for _, key := range cfg.v.AllKeys() {
			if key == "pre.key" {
				found = true
			}
		}

		if !found {
			t.Errorf("Expected [pre.key] to be available as config")
		}
	})

	t.Run("does not add a prefix to flag when prefix is added to ignoredPrefixes", func(t *testing.T) {
		e := &Engine{
			Cmd:       &cobra.Command{},
			ConfigKey: "pre",
			FlagSet:   pflag.NewFlagSet("dummy", pflag.ContinueOnError),
		}
		e.FlagSet.String("key", "", "")

		cfg := NewNutsGlobalConfig()
		cfg.IgnoredPrefixes = append(cfg.IgnoredPrefixes, "pre")
		cfg.RegisterFlags(e.Cmd, e)

		var found bool
		for _, key := range cfg.v.AllKeys() {
			println(key)
			if key == "key" {
				found = true
			}
		}

		if !found {
			t.Errorf("Expected [key] to be available as config")
		}
	})
}

func TestNutsGlobalConfig_LoadConfigFile(t *testing.T) {
	t.Run("Does not return error on missing file", func(t *testing.T) {
		cfg := NutsGlobalConfig{
			DefaultConfigFile: "non_existing.yaml",
			v:                 viper.New(),
		}
		cfg.Load(&cobra.Command{})

		if err := cfg.loadConfigFile(); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Returns error on incorrect file", func(t *testing.T) {
		cfg := NutsGlobalConfig{
			DefaultConfigFile: "../test/config/corrupt.yaml",
			v:                 viper.New(),
		}
		cfg.Load(&cobra.Command{})

		err := cfg.loadConfigFile()
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
			v:                 viper.New(),
		}
		cfg.Load(&cobra.Command{})

		err := cfg.loadConfigFile()
		if err != nil {
			t.Errorf("Expected no error, got [%v]", err.Error())
		}

		val := cfg.v.Get("key")
		if val != "value" {
			t.Errorf("Expected value to equals [value], got [%v]", val)
		}
	})
}

func TestNutsGlobalConfig_LoadAndUnmarshal(t *testing.T) {
	cfg := NewNutsGlobalConfig()
	cfg.Load(&cobra.Command{})

	t.Run("Adds configFile flag to Cmd", func(t *testing.T) {
		err := cfg.LoadAndUnmarshal(&cobra.Command{}, &struct{}{})
		if err != nil {
			t.Errorf("Expected no error, got [%v]", err.Error())
		}

		if !cfg.v.IsSet(configFileFlag) {
			t.Errorf("Expected %s to be set", configFileFlag)
		}
	})

	t.Run("Injects custom config into struct", func(t *testing.T) {
		s := struct {
			Key string
		}{
			"",
		}
		cfg.v.Set("key", "value")
		err := cfg.LoadAndUnmarshal(&cobra.Command{}, &s)
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
	cfg.Load(&cobra.Command{})

	t.Run("param is injected into engine without ConfigKey", func(t *testing.T) {
		c := struct {
			Key string
		}{}

		e := &Engine{
			Config:  &c,
			FlagSet: pflag.NewFlagSet("dummy", pflag.ContinueOnError),
		}
		e.FlagSet.String("key", "", "")

		cfg.v.Set("key", "value")

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
			Config:    &c,
			ConfigKey: "pre",
			FlagSet:   pflag.NewFlagSet("dummy", pflag.ContinueOnError),
		}
		e.FlagSet.String("key", "", "")
		cfg.RegisterFlags(e.Cmd, e)

		cfg.v.Set("pre.key", "value")

		if err := cfg.InjectIntoEngine(e); err != nil {
			t.Errorf("Expected no error, got [%v]", err.Error())
		}

		if c.Key != "value" {
			t.Errorf("Expected value to be injected into struct")
		}
	})

	t.Run("nested param is injected into engine without ConfigKey", func(t *testing.T) {
		c := struct {
			Nested struct{ Key string }
		}{}

		e := &Engine{
			Config:  &c,
			FlagSet: pflag.NewFlagSet("dummy", pflag.ContinueOnError),
		}
		e.FlagSet.String("nested.key", "", "")

		cfg.v.Set("nested.key", "value")

		if err := cfg.InjectIntoEngine(e); err != nil {
			t.Errorf("Expected no error, got [%v]", err.Error())
		}

		if c.Nested.Key != "value" {
			t.Errorf("Expected value to be injected into struct")
		}
	})

	t.Run("nested param is injected into engine with ConfigKey", func(t *testing.T) {
		c := struct {
			Nested struct{ Key string }
		}{}

		e := &Engine{
			Config:    &c,
			ConfigKey: "pre",
			FlagSet:   pflag.NewFlagSet("dummy", pflag.ContinueOnError),
		}
		e.FlagSet.String("nested.key", "", "")
		cfg.RegisterFlags(e.Cmd, e)

		cfg.v.Set("pre.nested.key", "value")

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
			Name:    "test",
			Config:  &c,
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
			Nested struct{ Key string }
		}{}

		e := &Engine{
			Config:  &c,
			FlagSet: pflag.NewFlagSet("dummy", pflag.ContinueOnError),
		}
		e.FlagSet.String("nested.key", "", "")

		cfg.v.Set("nested", "value")

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
			Config:  &c,
			FlagSet: pflag.NewFlagSet("dummy", pflag.ContinueOnError),
		}
		e.FlagSet.String("nested.key", "", "")

		cfg.v.Set("nested.key", "value")

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
