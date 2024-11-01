import json
import argparse
from typing import List
import requests
from sensor import Sensor
from sensortype import SensorType
from processingtype import ProcessingType

parser = argparse.ArgumentParser(
    description="Configure the server credentials")
parser.add_argument("-u", "--username", type=str,
                    help="Project Server username")
parser.add_argument("-p", "--password", type=str,
                    help="Project Server password")


def load_server_config():
    with open("config.json", "r") as file:
        config = json.load(file)

    args = parser.parse_args()
    username = args.username
    password = args.password

    if username is None or password is None:
        print("Username or password not configured as arguments, trying to parse them from config.json")
        username = config["server"].get("username")
        password = config["server"].get("password")

    srv_ip = config["server"].get("ip_addr")
    srv_port = config["server"].get("port")

    return srv_ip, srv_port, username, password


def login(srv_ip, srv_port, username, password) -> str:
    if not srv_ip or not srv_port:
        print("Server IP and port must be configured.")
        return None

    url = f"http://{srv_ip}:{srv_port}/api/v1/login"

    credentials = {
        'userName': username,
        'password': password
    }

    try:
        response = requests.post(url=url, json=credentials)
        response.raise_for_status()
        return response.json().get("auth_token", {}).get("token")
    except requests.exceptions.RequestException as e:
        print("Login failed:", e, response.json().get('error'))
        return None


def add_sensor_to_server(srv_ip: str, srv_port: str | int, token: str, sensor: Sensor) -> bool:
    url = f"http://{srv_ip}:{srv_port}/api/v1/sensor"

    headers = {
        "Authorization": f"Bearer {token}"
    }

    payload = {
        'name': sensor.name,
        'refresh_rate': sensor.refresh_rate,
        'uri': sensor.uri,
        'type': sensor.type.value
    }

    try:
        response = requests.post(url=url, headers=headers, json=payload)
        response.raise_for_status()
        return True
    except requests.exceptions.RequestException as e:
        print("Adding sensor failed:", e, response.json().get('error'))
        return False


def load_sensors_from_config() -> List[Sensor]:
    with open("config.json", "r") as file:
        config = json.load(file)

    sensors = []

    for sensor_data in config.get("sensors", []):
        name = sensor_data.get("name")
        ip_addr = sensor_data.get("ip_addr")
        port = sensor_data.get("port")
        refresh_rate = sensor_data.get("refresh_rate")

        processing_str = sensor_data.get("processing")
        processing = ProcessingType(processing_str) if processing_str else None
        number_of_hours = sensor_data.get("number_of_hours")

        params = sensor_data.get("params", {})

        sensor = Sensor(
            name=name,
            ip_addr=ip_addr,
            port=port,
            refresh_rate=refresh_rate,
            processing=processing,
            number_of_hours=number_of_hours,
            params=params
        )

        sensors.append(sensor)

    return sensors


if __name__ == '__main__':
    srv_ip, srv_port, username, password = load_server_config()
    token = login(srv_ip, srv_port, username, password)
    sensors = load_sensors_from_config()
    for sensor in sensors:
        add_sensor_to_server(srv_ip, srv_port, token, sensor)
