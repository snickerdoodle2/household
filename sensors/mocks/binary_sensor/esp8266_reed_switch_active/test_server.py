from flask import Flask, request

app = Flask(__name__)


@app.route('/activesensorupdate', methods=['POST'])
def active_sensor_update():
    # Get the parameters from the POST request
    value = request.form.get('value')
    sensor_id = request.form.get('sensor_id')

    # Print the values of the parameters
    print(f'Received value: {value}')
    print(f'Received sensor_id: {sensor_id}')

    # Return a simple response
    return "Data received", 200


if __name__ == '__main__':
    # Run the Flask app on port 5000
    app.run(host='10.0.0.55', port=5000)
