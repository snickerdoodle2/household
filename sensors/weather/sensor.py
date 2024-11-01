from processingtype import ProcessingType
from typing import Any
from sensortype import SensorType


class Sensor:
    def __init__(self, name: str, ip_addr: str, port: int, refresh_rate: int, processing: ProcessingType, number_of_hours: int, params: dict[str, Any]):
        self.name = name
        self.ip_addr = ip_addr
        self.port = port
        self.refresh_rate = refresh_rate
        self.processing = processing
        self.number_of_hours = number_of_hours
        self.params = params
        self.uri = f"{ip_addr}:{port}"
        self.type = SensorType.DECIMAL_SENSOR
