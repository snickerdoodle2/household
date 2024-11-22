import type { Result } from '@/types/result';
import { z } from 'zod';
import { authFetch } from './fetch';
import type { FetchFn } from '@/types/misc';
import { userSchema, type User } from '@/types/user';
import { sensorDetailsSchema, type SensorDetails } from '@/types/sensor';

export const getAllSensors = async (
    fetch: FetchFn
): Promise<Result<User[], string>> => {
    const res = await authFetch(`/api/v1/sensor/`, {}, fetch);
    const data = await res.json();
    if (!res.ok) {
        return {
            isError: true,
            error: data.error,
        };
    }

    const parsed = z.object({ data: userSchema.array() }).safeParse(data);
    if (!parsed.success) {
        return {
            isError: true,
            error: 'Error while parsing the data! (getAllUsers)',
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
