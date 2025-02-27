# GROWATT solar panels integration
Server communicating with GroWatt written in python. After logging in to GroWatt servers, server allows for retrieving data from your powerplant. 
Server is written in two versions - passive and active. Passive version returns last powerplant power reading. Active version sends last power reading, when new data is detected.

## Configuration
Configuration is done by growatt-sensor.properties file. You need to provide your plantId, username, and password for growatt server.

## Passive version endpoints
- `GET /value` -  mean powerplant power from last 5 minutes
- `GET /status` - status and sensor type
