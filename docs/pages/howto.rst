.. _nuts-crypto-howto:

Howto
=====

Using the library
-----------------

.. include:: ../../README.rst
    :start-after: .. inclusion-marker-for-contribution

Building the library
--------------------

For generating mocks:

.. code-block:: shell

   go get github.com/golang/mock/gomock
   go install github.com/golang/mock/mockgen

Then run

.. code-block:: shell

   mockgen -destination=mock/mock_oapi.go -package=mock github.com/deepmap/oapi-codegen/pkg/runtime EchoRouter
   mockgen -destination=mock/mock_echo.go -package=mock github.com/labstack/echo Context




