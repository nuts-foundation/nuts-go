nuts service executable
#######################

Nuts executable for Nuts service space. The idea behind this executable that it includes different 'engines'.
It can be configured through command line options to enable or disable an engine.
This will allow for a single process that runs all service space components, ideal for development.
For production a choice can be made for multiple instances of the same engine (by starting this executable multiple times), allowing for a more fine grained control and better scalability.
The executable exposes the REST (or other) services from the different engines. This also makes it easier to apply a particular security mechanism.

.. image:: https://travis-ci.org/nuts-foundation/nuts-go.svg?branch=master
    :target: https://travis-ci.org/nuts-foundation/nuts-go
    :alt: Build Status

.. image:: https://readthedocs.org/projects/nuts-go/badge/?version=latest
    :target: https://nuts-documentation.readthedocs.io/projects/nuts-go/en/latest/?badge=latest
    :alt: Documentation Status

.. image:: https://codecov.io/gh/nuts-foundation/nuts-go/branch/master/graph/badge.svg
    :target: https://codecov.io/gh/nuts-foundation/nuts-go

.. image:: https://api.codacy.com/project/badge/Grade/272258ac93e847b9b61c08d4144d0538
    :target: https://www.codacy.com/app/woutslakhorst/nuts-go

Dependencies
************

Go version => 1.13 is required.

Running tests
*************

Tests can be run by executing

.. code-block:: shell

    go test ./...

Building
********

just use ``go build``.

README
******

The readme is auto-generated from a template and uses the documentation to fill in the blanks.

.. code-block:: shell

    ./generate_readme.sh

This script uses ``rst_include`` which is installed as part of the dependencies for generating the documentation.

Documentation
*************

To generate the documentation, you'll need python3, sphinx and a bunch of other stuff. See :ref:`nuts-documentation-development-documentation`
The documentation can be build by running

.. code-block:: shell

    /docs $ make html

The resulting html will be available from ``docs/_build/html/index.html``

Configuration
*************

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
address                 --address               NUTS_ADDRESS            address and port server will be listening at.
configfile              --configfile            NUTS_CONFIGFILE         points to the location of the config file to be used.
verbosity               --verbosity             NUTS_VERBOSITY          Log level ("trace", "debug", "info", "warn", "error")
mode                    --mode                  NUTS_MODE               Mode the application will run in. When 'cli' it can be used to
                                                                        administer a remote Nuts node. When 'server' it will start a Nuts node.
                                                                        Defaults to 'server'.
=====================   ====================    =====================   ================================================================

