from flask import Flask, request, jsonify

app = Flask(__name__)


@app.route('/activesensorupdate', methods=['POST'])
def active_sensor_update():
    print("\n\n")
    print(request)
    print(request.get_json())
    print("\n\n")
    return jsonify({"message": "Data received"}), 200


if __name__ == '__main__':
    app.run(host='127.0.0.1', port=5000, debug=True)
