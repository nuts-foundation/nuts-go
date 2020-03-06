.. _nuts-go-config:

Nuts service config
###################

.. marker-for-readme

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
*************

The Nuts Global config adds the following options:

=====================   ====================    =====================   ================================================================
Name                    commandLine             Env                     Description
=====================   ====================    =====================   ================================================================
identity                --identity              NUTS_IDENTITY           Mandatory vendor identity (Vendor ID) of this node. It is the URN-encoded
                                                                        Chamber of Commerce registration no. of the vendor, e.g.:
                                                                        urn:oid:1.3.6.1.4.1.54851.4:12345678
address                 --address               NUTS_ADDRESS            address and port server will be listening at.
configfile              --configfile            NUTS_CONFIGFILE         points to the location of the config file to be used.
verbosity               --verbosity             NUTS_VERBOSITY          Log level ("trace", "debug", "info", "warn", "error")
mode                    --mode                  NUTS_MODE               Mode the application will run in. When 'cli' it can be used to
                                                                        administer a remote Nuts node. When 'server' it will start a Nuts node.
                                                                        Defaults to 'server'.
=====================   ====================    =====================   ================================================================
