import type { Result } from '@/types/result';
import {
    sensorDetailsSchema,
    sensorSchema,
    type Sensor,
    type SensorDetails,
} from '@/types/sensor';
import { z } from 'zod';
import { authFetch } from './fetch';
import type { FetchFn } from '@/types/misc';

export const SENSOR_URL = '/api/v1/sensor';

export const getAllSensors = async (
    fetch: FetchFn
): Promise<Result<Sensor[], string>> => {
    const res = await authFetch(SENSOR_URL, {}, fetch);
    const data = await res.json();
    if (!res.ok) {
        return {
            isError: true,
            error: data.error,
        };
    }

    const parsed = z.object({ data: sensorSchema.array() }).safeParse(data);
    if (!parsed.success) {
        return {
            isError: true,
            error: 'Error while parsing the data! (getAllSensors)',
        };
    }

    return {
        isError: false,
        data: parsed.data.data,
    };
};

export const getSensorDetails = async (
    id: string,
    fetch: FetchFn
): Promise<Result<SensorDetails, string>> => {
    const res = await authFetch(`${SENSOR_URL}/${id}`, {}, fetch);
    const data = await res.json();
    if (!res.ok) {
        return {
            isError: true,
            error: data.error,
        };
    }

    const parsed = z.object({ sensor: sensorDetailsSchema }).safeParse(data);
    if (!parsed.success) {
        return {
            isError: true,
            error: 'Error while parsing the data! (getSensorDetails)',
        };
    }

    return {
        isError: false,
        data: parsed.data.sensor,
    };
};
