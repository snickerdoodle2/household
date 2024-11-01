from sensortype import SensorType


class Sensor:
    def __init__(self, name: str, refresh_rate: int, uri: str, type: SensorType):
        self.name = name
        self.refresh_rate = refresh_rate
        self.uri = uri
        self.type = type
