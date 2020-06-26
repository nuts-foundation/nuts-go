.. _nuts-go-core-development:

Nuts service executable development
###################################

.. marker-for-readme

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

