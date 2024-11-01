import json
import argparse
import requests

parser = argparse.ArgumentParser(description="Configure the server credentials")
parser.add_argument("-u", "--username", type=str, help="Project Server username")
parser.add_argument("-p", "--password", type=str, help="Project Server password")

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

def login() -> str:
    srv_ip, srv_port, username, password = load_server_config()
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

if __name__ == '__main__':
    token = login()
    print(token)
