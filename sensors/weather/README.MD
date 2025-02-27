# Weather Data Server

The server is integrated with the Household project – after configuring and starting the server, sensors are automatically added to the Household application, the weather server starts running, and the standard operation of passive sensors begins.

The server uses the Open-Meteo API and allows access to all data found in the **_Hourly Weather Variables_** and **_Current Weather_** sections. Each sensor processes only one weather variable.

## Running the Server

Start the program from the command line using:

```
python weather-sensor-server.py
```

### Program Execution Parameters

- `-h, --help` - displays the help panel
- `-u, --username` - sets the username in _Household_
- `-p, --password` - sets the user password in _Household_
- `-c, --configpath` - sets the path to the configuration file

All parameters are optional. If the username or password is not provided, the program will attempt to read them from the configuration file. The default location for the configuration file is the server's root directory (from which the launch command is executed), and its name is `config.json`.

## Configuration

The server is configured through the `config.json` file. This file should define the parameters for the weather server, the Household application server, and the list of weather sensors in the correct format.

**Example configuration file snippet:**

```json
{
    "weather-server": {
        "host": "127.0.0.1",
        "port": 9001
    },
    "household-server": {
        "host": "172.30.227.16",
        "port": 8080,
        "username": "user",
        "password": "password"
    },
    "sensors": [
        {
            "name": "cloud-cover",
            "refresh_rate": 900,
            "processing": "current",
            "params": {
                "latitude": 50.3,
                "longitude": 19.28,
                "current": "cloud_cover"
            }
        },
        {
            "name": "rain sum next 24 hours",
            "refresh_rate": 3600,
            "processing": "sum_future",
            "number_of_hours": 24,
            "params": {
                "latitude": 50.3,
                "longitude": 19.28,
                "hourly": "rain",
                "forecast_days": 3
            }
        }
    ]
}
```

### Top-Level Structure

- `weather-server` - configuration for the host and port where the weather server runs
- `household-server` - configuration for connecting to the _Household_ system
- `sensors` - list of sensors, each with its respective attributes

### `weather-server` Object

- `host` - the host on which the weather server runs
- `port` - the port on which the weather server runs

### `household-server` Object

- `host` - the host on which the _Household_ system server runs
- `port` - the port on which the _Household_ system server runs
- `username` - [**optional**] the username for authentication with the _Household_ server, can also be provided as a program execution parameter
- `password` - [**optional**] the password for authentication with the _Household_ server, can also be provided as a program execution parameter

### `sensors` Array

The `sensors` array contains the configuration for each sensor. Each sensor is an object with the following fields:

- `name` - the sensor name
- `refresh-rate` - the data refresh rate in the _Household_ application (in seconds)
- `processing` - defines the type of data processing performed by the sensor. Allowed values:
  - `"sum_past"` - sum of past weather data values
  - `"sum_future"` - sum of forecasted weather data values
  - `"max_past"` - maximum of past weather data values
  - `"max_future"` - maximum of forecasted weather data values
  - `"min_past"` - minimum of past weather data values
  - `"min_future"` - minimum of forecasted weather data values
  - `"avg_past"` - average of past weather data values
  - `"avg_future"` - average of forecasted weather data values
  - `"current"` - current data

- `number_of_hours` - defines how many hours into the past or future (depending on the `processing` field) the data should be processed. If `processing` is `current`, this field is not required. It is important to ensure that the data obtained from _Open-Meteo_ has a sufficiently long timeframe to allow processing for the specified number of hours.

- `params` - the query parameters object sent to _Open-Meteo_. The easiest way to obtain this object is by selecting the desired data on [Open-Meteo's website](https://open-meteo.com/en/docs) and copying the `params` object from the _API Response_ > _Python_ section. Supported data includes _Hourly Weather Variables_ (for both _forecast days_ and _past days_ options) and _Current Weather_ data. It is crucial that this object contains only the single weather variable for which the sensor is responsible, as including additional variables may cause unexpected sensor behavior.

Each sensor is of type `decimal_sensor`.

## Requirements

- Python (tested on version 3.13.0)
- Required libraries:
  - `requests`
  - `flask`
  - `openmeteo_requests`
  - `numpy`
  - `pandas`

## Notes

- The `forecast_days` parameter in `params` specifies the number of forecasted days, including the current day. Using a value of `1` is not recommended, as it provides a forecast for only the current day. A value of `3` will return a forecast for the current day, the next day, and the day after that.

