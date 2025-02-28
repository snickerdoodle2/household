# Devices

Along with the system, software for several basic devices will be provided, as well as templates to facilitate programming custom devices. This will help beginner users enter the world of embedded systems and accelerate their development in using the application.

Sensors and effectors, which are external devices interacting with the system, can take any form. Their only requirement is to expose a REST API that is understood by the backend part of Household. This API can be exposed by Arduino boards managing any connected devices, "virtual" devices running as a server on a personal computer, in the cloud, or in any other way.

The system supports several basic types of devices that cover the full range of HouseHold applications. All of them can operate in passive mode (waiting for queries from the server – polling) or in active mode (sending information to the backend on their own, for example, when a sensor state change is detected).

Supported device types include:

- ``binary_switch`` - binary effector (e.g., light switch)
- ``binary_sensor`` - binary sensor (e.g., door open sensor)
- ``decimal_switch`` - decimal effector (e.g., dimmable light control)
- ``decimal_sensor`` - decimal sensor (e.g., thermometer)
- ``button`` - button (e.g., triggering a change in an automatic gate state)

### Passive Sensor API

Passive sensors must expose the following endpoints:

- ``GET /value`` - When queried via GET, the sensor returns a JSON response with the current measurement value. The value is always of type float.
- ``GET /status`` - Endpoint for verifying the device's status. If functioning correctly, a GET request will return a JSON indicating that the device is ONLINE and containing the device type.

### Passive Effector API

In addition to the passive sensor endpoints, passive effectors must include:

- ``POST /value`` - Upon receiving a request, the effector changes its value to the one provided in the request. An effector that is a button ignores the value and simulates button behavior (though a debounce time can also be sent as `value`).
- ``POST /toggle`` - Endpoint exposed only by binary effectors; upon receiving a request, it toggles the device state.

### Active Device API

Active devices, in addition to the endpoints provided by their passive counterparts, must also expose the 
- ``POST /init`` endpoint, allowing a handshake with the backend.

Active sensors perform a handshake initiated by the backend, during which they store information about the backend address, the backend endpoints where their measurements will be processed, and their unique identification token. This token must be sent with every measurement transmitted to the backend for it to be accepted and processed.

Sending measurement data to the backend involves making a POST request to the appropriate endpoint. After identifying the sensor, the processing of new values proceeds similarly to passive sensors.
