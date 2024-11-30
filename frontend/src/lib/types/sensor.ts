// created_at : "2024-04-15T16:04:49+02:00"
// id : "099317dc-edbf-40c6-9601-62f0b4921d75"
// name : "Drzwi wejściowe"
// refresh_rate : 5
// type : "binary_sensor"
// uri : "127.0.0.1:10001"
// version : 3

import { z } from 'zod';

export const sensorTypeSchema = z.enum([
    'binary_switch',
    'binary_sensor',
    'decimal_switch',
    'decimal_sensor',
    'button',
]);

export type SensorType = z.infer<typeof sensorTypeSchema>;

export const sensorSchema = z.object({
    id: z.string().uuid(),
    name: z.string(),
    type: sensorTypeSchema,
    active: z.boolean(),
});

export type Sensor = z.infer<typeof sensorSchema>;

export const newSensorSchema = z.object({
    name: z.string(),
    refresh_rate: z.number(),
    type: sensorTypeSchema,
    uri: z.string(),
    active: z.boolean(),
});

export type NewSensor = z.infer<typeof newSensorSchema>;

export const sensorDetailsSchema = newSensorSchema.merge(
    z.object({
        id: z.string().uuid(),
        created_at: z
            .string()
            .or(z.date())
            .transform((d) => new Date(d)),
    })
);

export type SensorDetails = z.infer<typeof sensorDetailsSchema>;
