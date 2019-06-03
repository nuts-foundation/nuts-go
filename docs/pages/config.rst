.. _nuts-go-config:

Config
======

The Nuts-go library contains some configuration logic which allows for usage of configFiles, Environment variables and commandLine params transparently.
If a Nuts engine is added as Engine it'll automatically work for the given engine. It is also possible for an engine to add the capabilities on a standalone basis.
This allows for testing from within a repo.

The parameters follow the following convention:
``$ nuts --parameter X`` is equal to ``$ NUTS_PARAMETER=X nuts`` is equal to ``parameter: X`` in a yaml file.

Or for this piece of yaml

.. code-block:: yaml

    nested:
        parameter: X

is equal to ``$ nuts --nested.parameter X`` is equal to ``$ NUTS_NESTED_PARAMETER=X nuts``

Config parameters for engines are prepended by the ``engine.ConfigKey`` by default (configurable):

.. code-block:: yaml

    engine:
        nested:
            parameter: X

is equal to ``$ nuts --engine.nested.parameter X`` is equal to ``$ NUTS_ENGINE_NESTED_PARAMETER=X nuts``


Added options
-------------

The Nuts Global config adds the following options:

=====================   ====================    =====================   ================================================================
Name                    commandLine             Env                     Description
=====================   ====================    =====================   ================================================================
configfile              --configfile            NUTS_CONFIGFILE         points to the location of the config file to be used.
More to be added
=====================   ====================    =====================   ================================================================


Usage
-----

The basic setup would be an engine with the following definition (fields not relevant for configuration are left out):

.. code-block:: go

    e := Engine {
        // The prefix used for any parameter in FlagSet, can be omitted or ignored in the NutsGlobalConfig
        ConfigKey string

        // Config is the pointer to a config struct. The struct can be used internally to read config values.
        Config interface{}

        // Configure checks if the combination of config parameters is allowed
        Configure func() error

        // FlagSet contains all engine-local configuration possibilities so they can be displayed through the help command
        // Load() also looks at this set of paramaters
        FlagSet *pflag.FlagSet
    }

Standalone
----------

To enable all Nuts commandLine options do

.. code-block:: go

    // engine instance
    var e = NewMyEngine()

    // the rootCmd
    var rootCmd = e.Cmd

    // a new global nuts config
    c := cfg.NewNutsGlobalConfig()

    // ignore any config prefixes for this Cmd since it is running standalone
    c.IgnoredPrefixes = append(c.IgnoredPrefixes, e.ConfigKey)

    // register all commandLine options added by this engine
    c.RegisterFlags(e)

    // load all config from parameters into global config
    if err := c.Load(); err != nil {
        panic(err)
    }

    // inject parameters from global config into config struct of engine
    if err := c.InjectIntoEngine(e); err != nil {
        panic(err)
    }

    // check configuration on engine
    if err := e.Configure(); err != nil {
        panic(err)
    }

    // execute comand
    rootCmd.Execute()
