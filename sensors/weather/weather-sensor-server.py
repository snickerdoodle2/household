import json
import argparse
import requests
from sensor import Sensor
from sensortype import SensorType

parser = argparse.ArgumentParser(
    description="Configure the server credentials")
parser.add_argument("-u", "--username", type=str,
                    help="Project Server username")
parser.add_argument("-p", "--password", type=str,
                    help="Project Server password")

with open("config.json", "r") as file:
    config = json.load(file)


def load_server_config():
    global srv_ip, srv_port, username, password

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
        print("Login failed:", e)
        return None


def add_sensor_to_server(srv_ip: str, srv_port: str | int, token: str, sensor: Sensor) -> bool:
    url = f"http://{srv_ip}:{srv_port}/api/v1/sensor"

    headers = {
        "Authorization": f"Bearer {token}"
    }

    payload = {
        'name': sensor.name,
        'refresh_rate': sensor.refresh_rate,
        'uri': "127.0.0.1:5005",
        'type': sensor.type.value
    }

    try:
        response = requests.post(url=url, headers=headers, json=payload)
        response.raise_for_status()
        return True
    except requests.exceptions.RequestException as e:
        print("Adding sensor failed:", e)
        return False


if __name__ == '__main__':
    srv_ip, srv_port, username, password = load_server_config()
    token = login(srv_ip, srv_port, username, password)
    print(token)
    test_sensor = Sensor("nazwa", 1200, "127.0.0.1:5555",
                         SensorType.BINARY_SENSOR)
    print(add_sensor_to_server(srv_ip, srv_port, token, test_sensor))
