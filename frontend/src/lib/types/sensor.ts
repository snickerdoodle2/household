// created_at : "2024-04-15T16:04:49+02:00"
// id : "099317dc-edbf-40c6-9601-62f0b4921d75"
// name : "Drzwi wejściowe"
// refresh_rate : 5
// type : "binary_sensor"
// uri : "127.0.0.1:10001"
// version : 3

import { z } from "zod";

export const SensorTypeEnum = z.enum([
    "binary_switch",
    "binary_sensor",
    "decimal_switch",
    "decimal_sensor",
    "button",
]);

export const Sensor = z.object({
    id: z.string().uuid(),
    name: z.string(),
    type: SensorTypeEnum,
    uri: z.string(),
    created_at: z.string().datetime({ offset: true }), // TODO: add day.js to this :)
    version: z.number(),
});

export type SensorType = z.infer<typeof Sensor>;
