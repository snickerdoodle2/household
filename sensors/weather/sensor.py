from processingtype import ProcessingType
from typing import Any
from sensortype import SensorType


class Sensor:
    def __init__(self, name: str, refresh_rate: int, processing: ProcessingType, number_of_hours: int, params: dict[str, Any]):
        self.name = name
        self.refresh_rate = refresh_rate
        self.processing = processing
        self.number_of_hours = number_of_hours
        self.params = params
        self.type = SensorType.DECIMAL_SENSOR
