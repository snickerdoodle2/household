# USAGE: flask run -h <your_ip>

from flask import Flask, request, jsonify

app = Flask(__name__)


@app.route('/activesensorupdate', methods=['POST'])
def active_sensor_update():
    # Get the JSON data from the POST request
    data = request.get_json()

    # Extract the 'value' and 'sensor_id' from the JSON data
    value = data.get('value')
    sensor_id = data.get('sensor_id')

    # Print the values of the parameters
    print(f'Received value: {value}')
    print(f'Received sensor_id: {sensor_id}')

    # Return a simple response
    return jsonify({"message": "Data received"}), 200


if __name__ == '__main__':
    # Run the Flask app on port 5000
    app.run(host='10.0.0.55', port=5000)
