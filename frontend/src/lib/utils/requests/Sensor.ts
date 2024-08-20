import { SERVER_URL } from '@/config/const';
import type { Result } from '@/types/Result.types';
import { SensorSchema, type Sensor } from '@/types/Sensor.types';
import { z } from 'zod';

export async function getAllSensors(
    fetch: (
        input: RequestInfo | URL,
        init?: RequestInit | undefined
    ) => Promise<Response>
): Promise<Result<Sensor[], string>> {
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
}

export async function getSensorData(sensorId: string) {
    // TODO: implement
    return Math.random();
}

// async function deleteDevice(id: Sensor['id']) {
//     console.log('Deleting sensor with id:', id);
//     try {
//         const response = await fetch(`${SERVER_URL}/api/v1/sensor/${id}`, {
//             method: 'DELETE',
//         });

//         if (response.ok) {
//         } else {
//             console.error('Error deleting sensor:', response.statusText);
//         }
//     } catch (error) {
//         console.error('Error deleting sensor:', error);
//     }

//     fetchSensors();
// }
