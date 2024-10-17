import requests
from datetime import datetime, timedelta
from jproperties import Properties
from flask import Flask


def login():
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

    print(login_data)

    try:
        response = requests.post(url, headers=headers, data=login_data)

        response.raise_for_status()

        cookies = response.cookies
        return cookies

    except requests.exceptions.RequestException as e:
        print(f"An error occurred during login: {e}")
        return None


def get_plant_detail(cookies):
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
        print(f"An error occurred during the plant detail request: {e}")
        return None


def get_most_recent_data(response_json):
    plant_data = response_json.get('back', {}).get('data', {})

    if not plant_data:
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
    if not cookies:
        cookies = login()
    response_json = get_plant_detail(cookies)
    if response_json:
        return str(get_most_recent_data(response_json))
    else:
        return "-1"


if __name__ == '__main__':
    api.run()
