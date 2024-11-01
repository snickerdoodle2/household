from enum import Enum, auto


class ProcessingType(Enum):
    SUM_PAST = "sum_past"
    SUM_FUTURE = "sum_future"
    MAX_PAST = "max_past"
    MAX_FUTURE = "max_future"
    MIN_PAST = "min_past"
    MIN_FUTURE = "min_future"
    AVG_PAST = "avg_past"
    AVG_FUTURE = "avg_future"
