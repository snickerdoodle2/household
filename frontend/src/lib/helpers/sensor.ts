import type { Result } from '@/types/result';
import { sensorSchema, type Sensor } from '@/types/sensor';
import { z } from 'zod';
import { authFetch } from './fetch';

export const getAllSensors = async (
    fetch: (
        input: RequestInfo | URL,
        init?: RequestInit | undefined
    ) => Promise<Response>
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
