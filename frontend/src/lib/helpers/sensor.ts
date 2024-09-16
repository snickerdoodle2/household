import type { Result } from '@/types/result';
import {
    sensorDetailsSchema,
    sensorSchema,
    type Sensor,
    type SensorDetails,
} from '@/types/sensor';
import { z } from 'zod';
import { authFetch } from './fetch';

type FetchFn = (
    input: RequestInfo | URL,
    init?: RequestInit | undefined
) => Promise<Response>;

export const getAllSensors = async (
    fetch: FetchFn
): Promise<Result<Sensor[], string>> => {
    const res = await authFetch(`/api/v1/sensor`, {}, fetch);
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
    const res = await authFetch(`/api/v1/sensor/${id}`, {}, fetch);
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
