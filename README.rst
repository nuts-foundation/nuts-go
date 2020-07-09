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

.. image:: https://api.codeclimate.com/v1/badges/2706f4616dbae18e8ea6/maintainability
   :target: https://codeclimate.com/github/nuts-foundation/nuts-go/maintainability
   :alt: Maintainability

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

The configuration options documentation is generated from the actual flags provided by the engines. When engines
are updated, this documentation should be regenerated to reflect any changes in provided flags. To regenerate the
configuration documentation run the following command from the project root:

.. code-block:: shell

    make update-docs

To build the documentation, you'll need python3, sphinx and a bunch of other stuff. See :ref:`nuts-documentation-development-documentation`
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


Options
*******

The following options can be configured:

========================================  ===================================================================================  ================================================================================================================================================================================
Key                                       Default                                                                              Description
========================================  ===================================================================================  ================================================================================================================================================================================
****
address                                   localhost:1323                                                                       Address and port the server will be listening to
configfile                                nuts.yaml                                                                            Nuts config file
identity                                                                                                                       Vendor identity for the node, mandatory when running in server mode. Must be in the format: urn:oid:1.3.6.1.4.1.54851.4:<number>
mode                                      server                                                                               Mode the application will run in. When 'cli' it can be used to administer a remote Nuts node. When 'server' it will start a Nuts node. Defaults to 'server'.
strictmode                                false                                                                                When set, insecure settings are forbidden.
verbosity                                 info                                                                                 Log level (trace, debug, info, warn, error)
**Auth**
auth.actingPartyCn                                                                                                             The acting party Common name used in contracts
auth.address                              localhost:1323                                                                       Interface and port for http server to bind to
auth.enableCORS                           false                                                                                Set if you want to allow CORS requests. This is useful when you want browsers to directly communicate with the nuts node.
auth.irmaConfigPath                                                                                                            path to IRMA config folder. If not set, a tmp folder is created.
auth.irmaSchemeManager                    pbdf                                                                                 The IRMA schemeManager to use for attributes. Can be either 'pbdf' or 'irma-demo'
auth.mode                                                                                                                      server or client, when client it does not start any services so that CLI commands can be used.
auth.publicUrl                                                                                                                 Public URL which can be reached by a users IRMA client
auth.skipAutoUpdateIrmaSchemas            false                                                                                set if you want to skip the auto download of the irma schemas every 60 minutes.
**ConsentBridgeClient**
cbridge.address                           http://localhost:8080                                                                API Address of the consent bridge
**ConsentStore**
cstore.address                            localhost:1323                                                                       Address of the server when in client mode
cstore.connectionstring                   \:memory:                                                                             Db connectionString
cstore.mode                                                                                                                    server or client, when client it uses the HttpClient
**Crypto**
crypto.fspath                             ./                                                                                   when file system is used as storage, this configures the path where keys are stored (default .)
crypto.keysize                            2048                                                                                 number of bits to use when creating new RSA keys
crypto.storage                            fs                                                                                   storage to use, 'fs' for file system (default)
**Events octopus**
events.autoRecover                        false                                                                                Republish unfinished events at startup
events.connectionstring                   file::memory:?cache=shared                                                           db connection string for event store
events.incrementalBackoff                 8                                                                                    Incremental backoff per retry queue, queue 0 retries after 1 second, queue 1 after {incrementalBackoff} * {previousDelay}
events.maxRetryCount                      5                                                                                    Max number of retries for events before giving up (only for recoverable errors
events.natsPort                           4222                                                                                 Port for Nats to bind on
events.purgeCompleted                     false                                                                                Purge completed events at startup
events.retryInterval                      60                                                                                   Retry delay in seconds for reconnecting
**Network**
network.address                                                                                                                Interface and port for http server to bind to, defaults to global Nuts address.
network.bootstrapNodes                                                                                                         Space-separated list of bootstrap nodes (`<host>:<port>`) which the node initially connect to.
network.certFile                                                                                                               PEM file containing the certificate this node will identify itself with to other nodes. If not set, the Nuts node will attempt to load a TLS certificate from the crypto module.
network.certKeyFile                                                                                                            PEM file containing the key belonging to this node's certificate. If not set, the Nuts node will attempt to load a TLS certificate from the crypto module.
network.grpcAddr                          \:5555                                                                                Local address for gRPC to listen on.
network.mode                                                                                                                   server or client, when client it uses the HttpClient
network.nodeID                                                                                                                 Instance ID of this node under which the public address is registered on the nodelist. If not set, the Nuts node's identity will be used.
network.publicAddr                                                                                                             Public address (of this node) other nodes can use to connect to it. If set, it is registered on the nodelist.
**Registry**
registry.address                          localhost:1323                                                                       Interface and port for http server to bind to, default: localhost:1323
registry.clientTimeout                    10                                                                                   Time-out for the client in seconds (e.g. when using the CLI), default: 10
registry.datadir                          ./data                                                                               Location of data files, default: ./data
registry.mode                             server                                                                               server or client, when client it uses the HttpClient, default: server
registry.organisationCertificateValidity  365                                                                                  Number of days organisation certificates are valid, default: 365
registry.syncAddress                      https://codeload.github.com/nuts-foundation/nuts-registry-development/tar.gz/master  The remote url to download the latest registry data from, default: https://codeload.github.com/nuts-foundation/nuts-registry-development/tar.gz/master
registry.syncInterval                     30                                                                                   The interval in minutes between looking for updated registry files on github, default: 30
registry.syncMode                         fs                                                                                   The method for updating the data, 'fs' for a filesystem watch or 'github' for a periodic download, default: fs
registry.vendorCACertificateValidity      1095                                                                                 Number of days vendor CA certificates are valid, default: 1095
**Validation**
fhir.schemapath                                                                                                                location of json schema, default nested Asset
========================================  ===================================================================================  ================================================================================================================================================================================

