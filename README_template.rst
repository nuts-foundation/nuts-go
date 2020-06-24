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

.. include:: docs/pages/development/nuts-go.rst
    :start-after: .. marker-for-readme

Configuration
*************

.. include:: docs/pages/configuration/nuts-go.rst
    :start-after: .. marker-for-readme
