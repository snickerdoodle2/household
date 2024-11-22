from enum import Enum


class SensorType(Enum):
    BINARY_SWITCH = "binary_switch"
    BINARY_SENSOR = "binary_sensor"
    DECIMAL_SWITCH = "decimal_switch"
    DECIMAL_SENSOR = "decimal_sensor"
    BUTTON = "button"
