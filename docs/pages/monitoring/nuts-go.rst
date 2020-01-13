.. _nuts-go-core-monitoring:

Nuts service executable monitoring
##################################

Basic service health
********************

A status endpoint is provided to check if the service is running and if the web server has been started.
The endpoint is available over http so it can be used by a wide range of health checking services.
It does not provide any information on the individual engines running as part of the executable.
The main goal of the service is to give a YES/NO answer for if the service is running?

.. code-block::

    GET /status

It'll return an "OK" response and a 200 status code.