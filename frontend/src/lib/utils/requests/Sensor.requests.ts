import { SERVER_URL } from '@/config/const';
import type { Result } from '@/types/Result.types';
import { SensorSchema, type Sensor, type SensorData } from '@/types/Sensor.types';
import { z } from 'zod';

export async function getAllSensors(
    fetch: (input: RequestInfo | URL, init?: RequestInit | undefined) => Promise<Response>
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

export async function removeSensor(id: Sensor['id']): Promise<Result<string, string>> {
    try {
        const response = await fetch(`${SERVER_URL}/api/v1/sensor/${id}`, {
            method: 'DELETE',
        });

        if (response.ok) {
            return {
                isError: false,
                data: 'Sensor removed successfully',
            };
        } else {
            return {
                isError: true,
                error: `Error: ${response.statusText}`,
            };
        }
    } catch (error) {
        return {
            isError: true,
            error: `Network error. Please try again later. ${error}`,
        };
    }
}

export async function addSensor({
    name,
    uri,
    type,
    refresh_rate,
}: SensorData): Promise<Result<SensorData, string>> {
    try {
        const response = await fetch(`${SERVER_URL}/api/v1/sensor`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name,
                uri,
                type,
                refresh_rate,
            }),
        });

        if (!response.ok) {
            // TODO: error json handling
            const errorData = await response.json();
            return {
                isError: true,
                error: `Error: ${errorData.error}`,
            };
        } else {
            // TODO: nice pop-up window instead of alert
            const responseData = await response.json();
            return {
                isError: false,
                data: responseData,
            };
        }
    } catch (error) {
        return {
            isError: true,
            error: 'Network error. Please try again later.',
        };
    }
}

export async function modifySensor(
    id: Sensor['id'],
    { name, uri, type, refresh_rate }: SensorData
): Promise<Result<SensorData, string>> {
    try {
        const response = await fetch(`${SERVER_URL}/api/v1/sensor/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name, uri, type, refresh_rate }),
        });

        if (!response.ok) {
            // TODO: error json handling
            const errorData = await response.json();
            return {
                isError: true,
                error: `Error: ${errorData.error}`,
            };
        } else {
            // TODO: nice pop-up window instead of alert
            const responseData = await response.json();
            return {
                isError: false,
                data: responseData,
            };
        }
    } catch (error) {
        return {
            isError: true,
            error: `Network error. Please try again later. ${error}`,
        };
    }
}
