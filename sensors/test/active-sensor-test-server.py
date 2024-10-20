from flask import Flask, request, jsonify

app = Flask(__name__)


@app.route('/activesensorupdate', methods=['POST'])
def active_sensor_update():
    data = request.get_json()

    value = data.get('value')
    sensor_id = data.get('sensor_id')

    print(f'Received value: {value}')
    print(f'Received sensor_id: {sensor_id}')

    return jsonify({"message": "Data received"}), 200


if __name__ == '__main__':
    app.run(host='127.0.0.1', port=5000)
