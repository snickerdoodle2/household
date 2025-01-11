import requests
from datetime import datetime, timedelta
from jproperties import Properties
from flask import Flask, jsonify

plantid = None


def login():
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
        print("plant data error")
        return 0

    now = datetime.now()

    for timestamp, value in plant_data.items():
        time_obj = datetime.strptime(timestamp, '%Y-%m-%d %H:%M')

        if timedelta(minutes=0) <= now - time_obj < timedelta(minutes=5):
            return value
    return 0


cookies = None

api = Flask(__name__)


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

    print(type(get_most_recent_value(response_json)))
    return jsonify(value=float(get_most_recent_value(response_json)))


@api.route('/status', methods=['GET'])
def get_status():
    response = jsonify(status="online",
                       type="decimal_sensor")
    return response, 200


if __name__ == '__main__':
    api.run(host="0.0.0.0")
