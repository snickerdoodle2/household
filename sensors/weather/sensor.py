from processingtype import ProcessingType
from sensortype import SensorType
from typing import Any
import openmeteo_requests
from openmeteo_requests.Client import WeatherApiResponse
import pandas as pd
import openmeteo_requests
import numpy as np
from datetime import datetime, timedelta, timezone


class Sensor:
    def __init__(self, name: str, refresh_rate: int, processing: ProcessingType, number_of_hours: int, params: dict[str, Any]):
        self.name = name
        self.refresh_rate = refresh_rate
        self.processing = processing
        self.number_of_hours = -number_of_hours if processing in {
            ProcessingType.SUM_PAST, ProcessingType.MAX_PAST, ProcessingType.MIN_PAST, ProcessingType.AVG_PAST} else number_of_hours
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
        return False

    def get_value(self):
        weather_data = self.__get_weather_data()

        if self.processing != ProcessingType.CURRENT:
            current_time = datetime.now(timezone.utc)
            start_time = min(current_time, current_time +
                             timedelta(hours=self.number_of_hours))
            end_time = max(current_time, current_time +
                           timedelta(hours=self.number_of_hours))

            values = self.__get_values_in_date_range(
                weather_data, start_time, end_time)

            return self.__process_values(values)
        else:
            return weather_data.Current().Variables(0).Value()

    def __get_weather_data(self):
        openmeteo = openmeteo_requests.Client()
        url = "https://api.open-meteo.com/v1/forecast"

        responses = openmeteo.weather_api(url, params=self.params)
        return responses[0]

    def __get_values_in_date_range(self, response: WeatherApiResponse, start_time: datetime, end_time: datetime) -> pd.DataFrame:
        hourly = response.Hourly()
        hourly_value = hourly.Variables(0).ValuesAsNumpy()

        hourly_data = {"date": pd.date_range(
            start=pd.to_datetime(hourly.Time(), unit="s", utc=True),
            end=pd.to_datetime(hourly.TimeEnd(), unit="s", utc=True),
            freq=pd.Timedelta(seconds=hourly.Interval()),
            inclusive="left"
        )}
        hourly_data["value"] = hourly_value

        hourly_dataframe = pd.DataFrame(data=hourly_data)

        hourly_dataframe['date'] = pd.to_datetime(hourly_dataframe['date'])

        return hourly_dataframe[(hourly_dataframe['date'] >= start_time) & (
            hourly_dataframe['date'] <= end_time)]

    def __process_values(self, df: pd.DataFrame) -> float:
        match(self.processing):
            case(ProcessingType.SUM_FUTURE | ProcessingType.SUM_PAST):
                return df["value"].sum()
            case(ProcessingType.MAX_PAST | ProcessingType.MAX_FUTURE):
                return df["value"].max()
            case(ProcessingType.MIN_FUTURE | ProcessingType.MAX_FUTURE):
                return df["value"].min()
            case(ProcessingType.AVG_FUTURE | ProcessingType.AVG_FUTURE):
                return df["value"].mean()
