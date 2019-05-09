nuts-go
===========

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

.. image:: https://api.codacy.com/project/badge/Grade/d17595243fd7424fbe743da68c8dcbfc
    :target: https://www.codacy.com/app/woutslakhorst/nuts-go

.. inclusion-marker-for-contribution

Installation
------------

.. code-block:: shell

   go get github.com/nuts-foundation/nuts-go
   make

Configuration
-------------

The lib is configured using `Viper <https://github.com/spf13/viper>`_ and `Cobra <https://github.com/spf13/cobra>`_.
Run:

.. code-block:: shell

   nuts --help

for listing the different options
