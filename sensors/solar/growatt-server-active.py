import json
import threading
import requests
from datetime import datetime, timedelta
import time
from jproperties import Properties
from flask import Flask, jsonify, request
import logging

api = Flask(__name__)

cookies = None
last_response_json = None
is_running = True
plantid = None

server_uri = None
server_measurements_endpoint = None
server_init_ack_endpoint = None
id_token = None

init_event = threading.Event()
config_lock = threading.Lock()

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


def login():
    logger.info("logging in")
    configs = Properties()
    with open('growatt-sensor.properties', 'rb') as config_file:
        configs.load(config_file)

    global plantid
    plantid = configs.get("PLANTID").data

    url = 'http://server.growatt.com/LoginAPI.do'
    headers = {
        'Content-Type': 'application/x-www-form-urlencoded',
        'User-Agent': 'Mozilla/5.0 (compatible; Linux; Python Requests)'
    }
    login_data = {
        'userName': (configs.get("USERNAME").data),
        'password': (configs.get("PASSWORD").data)
    }

    try:
        response = requests.post(url, headers=headers, data=login_data)
        response.raise_for_status()
        return response.cookies
    except requests.exceptions.RequestException as e:
        return None


def get_plant_detail(cookies):
    """This function returns an unordered list of average production power in 5 minute timeframes
    starting from date's sunrise ending in current hour or given date's sunset, the most recent
    timeframe shows up with about one minute delay
    """
    plant_url = 'http://server.growatt.com/newPlantDetailAPI.do'

    params = {
        'plantId': plantid,
        'type': '1',
        'date': datetime.now().strftime('%Y-%m-%d')
    }

    headers = {
        'Content-Type': 'application/x-www-form-urlencoded',
        'User-Agent': 'Mozilla/5.0 (compatible; Linux; Python Requests)'
    }

    try:
        response = requests.get(plant_url, params=params,
                                headers=headers, cookies=cookies)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        return None


def get_most_recent_value(response_json):
    plant_data = response_json.get('back', {}).get('data', {})

    if not plant_data:
        return 0

    now = datetime.now()

    for timestamp, value in plant_data.items():
        time_obj = datetime.strptime(timestamp, '%Y-%m-%d %H:%M')
        if timedelta(minutes=0) <= now - time_obj < timedelta(minutes=5):
            return value
    return 0


def send_init_ack(input_data):
    """Send initialization acknowledgment to the server"""
    client = requests.Session()
    client.timeout = 5

    with config_lock:
        url = server_uri
        init_ack_endpoint = server_init_ack_endpoint

    if not url.startswith(('http://', 'https://')):
        url = f"http://{url}{init_ack_endpoint}"

    try:
        response = client.post(
            url,
            json=input_data,
            headers={'Content-Type': 'application/json'}
        )
        logger.info(f"Sending init ack to {url}")

        if response.status_code < 300:
            return None
        return f"Init response code: {response.status_code}"
    except requests.exceptions.RequestException as e:
        logger.error(f"Error sending init ack: {e}")
        return str(e)


def send_measurement(value):
    """Send measurement to the configured server"""
    with config_lock:
        url = server_uri
        measurements_endpoint = server_measurements_endpoint
        local_id_token = id_token

    if not url.startswith(('http://', 'https://')):
        url = f"http://{url}{measurements_endpoint}"

    measurement = {
        "message-type": "measurement",
        "sensor-type": "decimal_sensor",
        "value": value,
        "id-token": local_id_token
    }

    try:
        response = requests.post(
            url,
            json=measurement,
            headers={'Content-Type': 'application/json'},
            timeout=10
        )
        logger.info(f"Sending measurement to: {url}")

        if response.status_code >= 300:
            logger.error(f"Server returned error status: {
                         response.status_code}")
    except requests.exceptions.RequestException as e:
        logger.error(f"Error sending measurement: {e}")


def fetch_data_periodically(app):
    """Background thread for fetching and sending data"""
    global cookies, last_response_json, is_running

    logger.info("Waiting for initialization...")
    init_event.wait()
    logger.info("Initialized - starting measurements")

    with app.app_context():
        while is_running:
            now = datetime.now()
            minutes = now.minute

            if minutes % 5 == 2:
                if not cookies or 'JSESSIONID' not in cookies:
                    cookies = login()
                    if not cookies or 'JSESSIONID' not in cookies:
                        logger.error("GroWatt login error")
                        continue

                response_json = get_plant_detail(cookies)
                if not response_json:
                    logger.error("Error fetching plant details")
                    continue

                if response_json != last_response_json:
                    last_response_json = response_json
                    value = get_most_recent_value(last_response_json)
                    logger.info(f"{now} Data updated, sending measurement")
                    send_measurement(value)

            time.sleep(60)


@api.route('/init', methods=['POST'])
def init_handler():
    """Handle initialization requests"""
    if not request.is_json:
        return "Invalid JSON", 400

    input_data = request.get_json()

    with config_lock:
        global server_uri, server_measurements_endpoint, server_init_ack_endpoint, id_token
        server_uri = input_data.get('server-uri')
        server_measurements_endpoint = input_data.get('measurements-endpoint')
        server_init_ack_endpoint = input_data.get('init-ack-endpoint')
        id_token = input_data.get('id-token')

    logger.info(f"Init request received: {input_data}")

    error = send_init_ack(input_data)
    if error:
        return f"Error sending init ack: {error}", 500

    init_event.set()

    return "", 200


@api.route('/status', methods=['GET'])
def get_status():
    return jsonify(status="online", type="decimal_sensor"), 200


@api.route('/value', methods=['GET'])
def get_value():
    global cookies

    if not cookies or 'JSESSIONID' not in cookies:
        cookies = login()
        if not cookies or 'JSESSIONID' not in cookies:
            logger.error("login error")
            return "GroWatt login error", 500

    response_json = get_plant_detail(cookies)
    if not response_json:
        logger.error("plant details error")
        return "Get plant detail error", 500

    return jsonify(value=get_most_recent_value(response_json))


if __name__ == '__main__':
    data_fetch_thread = threading.Thread(
        target=fetch_data_periodically, args=(api,))
    data_fetch_thread.daemon = True
    data_fetch_thread.start()

    try:
        api.run(debug=True, port=5023)
    finally:
        is_running = False
