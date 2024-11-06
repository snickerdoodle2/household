import json
import argparse
from typing import Dict, List
import requests
from sensor import Sensor
from processingtype import ProcessingType
from flask import Flask, jsonify


app = Flask(__name__)
sensors: Dict[str, Sensor] = {}
self_host = None
self_port = None
logger = app.logger
logger.setLevel("INFO")


parser = argparse.ArgumentParser(
    description="Configure the Weather Data Server")
parser.add_argument("-u", "--username", type=str,
                    help="Household Server username")
parser.add_argument("-p", "--password", type=str,
                    help="Household Server password")
parser.add_argument("-c", "--configpath", type=str,
                    help = "Path to config.json file")

args = parser.parse_args()
config_path = args.configpath if args.configpath is not None else "config.json"

def get_household_server_config():
    with open(config_path, "r") as file:
        config = json.load(file)

    username = args.username
    password = args.password

    if username is None or password is None:
        logger.info("Username or password not configured as arguments, trying to parse them from config.json")
        username = config["household-server"].get("username")
        password = config["household-server"].get("password")

    srv_host = config["household-server"].get("host")
    srv_port = config["household-server"].get("port")

    return srv_host, srv_port, username, password


def load_weather_server_config():
    with open(config_path, "r") as file:
        config = json.load(file)

    global self_host
    global self_port

    self_host = config["weather-server"].get("host")
    self_port = config["weather-server"].get("port")


def login(household_host, household_port, username, password) -> str:
    if not household_host or not household_port:
        logger.error("Server IP and port must be configured.")
        return None

    url = f"http://{household_host}:{household_port}/api/v1/login"

    credentials = {
        'userName': username,
        'password': password
    }

    try:
        response = requests.post(url=url, json=credentials)
        response.raise_for_status()
        return response.json().get("auth_token", {}).get("token")
    except requests.exceptions.ConnectionError as e:
        logger.error("Login failed: Unable to connect to the Household server.")
        exit(0)
    except requests.exceptions.RequestException as e:
        error_message = (
            response.json().get('error') if 'response' in locals() and response else str(e)
        )
        logger.error("Login to Household failed:", e, error_message)
        exit(0)


def add_sensor_to_household(household_ip: str, household_port: str | int, auth_token: str, sensor: Sensor) -> bool:
    url = f"http://{household_ip}:{household_port}/api/v1/sensor"
    global self_host, self_port

    headers = {
        "Authorization": f"Bearer {auth_token}"
    }

    payload = {
        'name': sensor.name,
        'refresh_rate': sensor.refresh_rate,
        'uri': f"{self_host}:{self_port}/{sensor.name}",
        'type': sensor.type.value
    }

    try:
        response = requests.post(url=url, headers=headers, json=payload)
        response.raise_for_status()
        return True
    except requests.exceptions.RequestException as e:
        error_message = response.json().get('error', {})

        if error_message.get('uri') == 'a sensor with this URI already exists':
            logger.warning(f"Sensor with URI {
                  payload['uri']} already exists. Continuing without adding duplicate.")
            return True

        logger.error("Adding sensor to Household failed:", e, error_message)
        return False


def get_sensors_from_config() -> List[Sensor]:
    with open(config_path, "r") as file:
        config = json.load(file)

    config_sensors = []

    for sensor_data in config.get("sensors", []):
        name = sensor_data.get("name")
        refresh_rate = sensor_data.get("refresh_rate")
        processing = ProcessingType(sensor_data.get("processing"))
        number_of_hours = sensor_data.get("number_of_hours")
        params = sensor_data.get("params", {})

        sensor = Sensor(
            name=name,
            refresh_rate=refresh_rate,
            processing=processing,
            number_of_hours=number_of_hours,
            params=params
        )
        if sensor.validate_params():
            config_sensors.append(sensor)
        else:
            logger.warning(
                f"sensor {sensor.name} failed param validation and will not be processed")

    return config_sensors


def initialize_sensors():
    global sensors
    household_ip, household_port, username, password = get_household_server_config()
    token = login(household_ip, household_port, username, password)
    sensor_list = get_sensors_from_config()
    for sensor in sensor_list:
        if add_sensor_to_household(household_ip, household_port, token, sensor):
            sensors[sensor.name] = sensor


@app.route("/<sensor_name>/status", methods=["GET"])
def get_sensor_status(sensor_name):
    global sensors
    sensor = sensors.get(sensor_name)
    if sensor:
        response = jsonify(status="online",
                           type="decimal_sensor")
        return response, 200
    else:
        return jsonify({"error": "Sensor not found"}), 404


@app.route("/<sensor_name>/value", methods=["GET"])
def get_sensor_value(sensor_name):
    global sensors
    sensor = sensors.get(sensor_name)
    if sensor:
        value = float(sensor.get_value())
        response = jsonify(value=value)
        return response, 200
    else:
        return jsonify({"error": "Sensor not found"}), 404


if __name__ == '__main__':
    load_weather_server_config()
    initialize_sensors()
    app.run(host=self_host, port=self_port)
