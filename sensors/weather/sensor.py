from processingtype import ProcessingType
from sensortype import SensorType
from typing import Any


class Sensor:
    def __init__(self, name: str, refresh_rate: int, processing: ProcessingType, number_of_hours: int, params: dict[str, Any]):
        self.name = name
        self.refresh_rate = refresh_rate
        self.processing = processing
        self.number_of_hours = number_of_hours
        self.params = params
        self.type = SensorType.DECIMAL_SENSOR


    def validate_params(self) -> bool:
        past_days = self.params.get("past_days")
        forecast_days = self.params.get("forecast_days")
        hourly = self.params.get("hourly")
        current = self.params.get("current")

        if self.processing in {ProcessingType.SUM_PAST, ProcessingType.MAX_PAST, ProcessingType.MIN_PAST, ProcessingType.AVG_PAST}:
            return past_days is not None and past_days * 24 >= self.number_of_hours

        elif self.processing in {ProcessingType.SUM_FUTURE, ProcessingType.MAX_FUTURE, ProcessingType.MIN_FUTURE, ProcessingType.AVG_FUTURE}:
            return (forecast_days is not None and forecast_days * 24 >= self.number_of_hours) or \
                (hourly is not None and hourly * 24 >= self.number_of_hours)

        elif self.processing == ProcessingType.CURRENT:
            return current is not None
        print(34)
        return False
