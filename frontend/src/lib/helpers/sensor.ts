import type { Result } from '@/types/result';
import { Sensor, type SensorType } from '@/types/sensor';
import { z } from 'zod';
import { authFetch } from './fetch';

export const getAllSensors = async (
    fetch: (
        input: RequestInfo | URL,
        init?: RequestInit | undefined
    ) => Promise<Response>
): Promise<Result<SensorType[], string>> => {
    const res = await authFetch(`/api/v1/sensor`, {}, fetch);
    const data = await res.json();
    if (!res.ok) {
        return {
            isError: true,
            error: data.error,
        };
    }

    const parsed = z.object({ data: Sensor.array() }).safeParse(data);
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
