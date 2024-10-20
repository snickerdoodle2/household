import json
import threading
import requests
from datetime import datetime, timedelta
import time
from jproperties import Properties
from flask import Flask, jsonify

api = Flask(__name__)

cookies = None
last_response_json = None

is_running = True

# test server endpoint - to be changed
POST_URL = "http://127.0.0.1:5000/activesensorupdate"


def login():
    print("logging in")
    configs = Properties()
    with open('credentials.properties', 'rb') as config_file:
        configs.load(config_file)

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

        cookies = response.cookies
        return cookies

    except requests.exceptions.RequestException as e:
        return None


def get_plant_detail(cookies):
    """this function returns an undordered list of average production power in 5 minute timeframes
    starting from date's sunrise ending in current hour or givens date sunset, the most recent
    timeframe shows up with about one minute delay
    """

    plant_url = 'http://server.growatt.com/newPlantDetailAPI.do'

    params = {
        'plantId': '418844',
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


def fetch_data_periodically(app):
    """
    This function runs in a background thread, checking for new data every 5 minutes when minutes % 5 == 2.
    If new data is found, it sends a POST request with the data.
    """
    global cookies, last_response_json, is_running

    with app.app_context():
        while is_running:
            now = datetime.now()
            minutes = now.minute

            # on growatt server data is updated every 5 minutes, but with various delay
            # checking delay of 2-3 minutes should ensure that new data is always available
            if minutes % 5 == 2:
                if not cookies or 'JSESSIONID' not in cookies:
                    cookies = login()
                    if not cookies or 'JSESSIONID' not in cookies:
                        print("GroWatt login error")
                        continue

                response_json = get_plant_detail(cookies)
                if not response_json:
                    print("Error fetching plant details")
                    continue

                if response_json != last_response_json:
                    last_response_json = response_json
                    print(now, "Data updated, sending POST request")
                    send_post_request(get_most_recent_value(last_response_json))

            time.sleep(60)


def send_post_request(value):
    """
    Sends a POST request with the data to the specified POST_URL.
    """
    try:
        headers = {'Content-Type': 'application/json'}
        data = json.dumps({"sensor_id": "growatt-5min", "value": value})  # Create JSON payload
        response = requests.post(POST_URL, data=data, headers=headers)
        if response.status_code == 200:
            print("Data sent successfully")
        else:
            print(f"Failed to send data: {response.status_code}")
    except Exception as e:
        # print(f"Error sending POST request: {e}")
        pass


@api.route('/value', methods=['GET'])
def get_value():
    global cookies

    if not cookies or 'JSESSIONID' not in cookies:
        cookies = login()
        if not cookies or 'JSESSIONID' not in cookies:
            print("login error")
            return "GroWatt login error", 500

    response_json = get_plant_detail(cookies)

    if not response_json:
        print("plant details error")
        return "Get plant detail error", 500

    return jsonify(value=get_most_recent_value(response_json))


@api.route('/status', methods=['GET'])
def get_status():
    response = jsonify(status="online",
                       type="decimal_sensor")
    return response, 200


if __name__ == '__main__':
    login()
    data_fetch_thread = threading.Thread(target=fetch_data_periodically, args=(api,))
    data_fetch_thread.daemon = True
    data_fetch_thread.start()

    try:
        api.run(debug=True, port=5023)
    finally:
        is_running = False
