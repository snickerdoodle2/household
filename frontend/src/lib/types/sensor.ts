// created_at : "2024-04-15T16:04:49+02:00"
// id : "099317dc-edbf-40c6-9601-62f0b4921d75"
// name : "Drzwi wejściowe"
// refresh_rate : 5
// type : "binary_sensor"
// uri : "127.0.0.1:10001"
// version : 3

import { z } from 'zod';

export enum SensorType {
    BINARY_SWITCH = 'binary_switch',
    BINARY_SENSOR = 'binary_sensor',
    DECIMAL_SWITCH = 'decimal_switch',
    DECIMAL_SENSOR = 'decimal_sensor',
    BUTTON = 'button',
}

export const SensorTypeSchema = z.nativeEnum(SensorType);

export const SensorSchema = z.object({
    id: z.string().uuid(),
    name: z.string(),
    type: SensorTypeSchema,
    uri: z.string(),
    created_at: z.string().datetime({ offset: true }), // TODO: add day.js to this :)
    version: z.number(),
    refresh_rate: z.number(),
});

export type Sensor = z.infer<typeof SensorSchema>;
