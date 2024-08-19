import { SERVER_URL } from '@/config/const';
import type { Result } from '@/types/Result.types';
import { SensorSchema, type Sensor } from '@/types/Sensor.types';
import { z } from 'zod';

export const getAllSensors = async (
    fetch: (
        input: RequestInfo | URL,
        init?: RequestInit | undefined
    ) => Promise<Response>
): Promise<Result<Sensor[], string>> => {
    const res = await fetch(`${SERVER_URL}/api/v1/sensor`);
    const data = await res.json();
    if (!res.ok) {
        return {
            isError: true,
            error: data.error,
        };
    }

    const parsed = z.object({ data: SensorSchema.array() }).safeParse(data);
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

async function fetchSensors() {
    const response = await getAllSensors(fetch);
    if (response.isError) {
        console.error('Failed to fetch sensors!', response.error);
        return;
    }
    sensorStore.set(response.data);
    loading = false;
}

async function deleteDevice(id: Sensor['id']) {
    console.log('Deleting sensor with id:', id);
    try {
        const response = await fetch(`${SERVER_URL}/api/v1/sensor/${id}`, {
            method: 'DELETE',
        });

        if (response.ok) {
        } else {
            console.error('Error deleting sensor:', response.statusText);
        }
    } catch (error) {
        console.error('Error deleting sensor:', error);
    }

    fetchSensors();
}
