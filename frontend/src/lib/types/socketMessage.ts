import { z } from 'zod';
import { notificationSchema } from './notification';

export const sensorSocketMessageSchema = z.object({
    values: z.number().array(),
    status: z.enum(['ONLINE', 'OFFLINE']),
});

export type SensorSocketMessage = z.infer<typeof sensorSocketMessageSchema>;

export const notificationSocketMessageSchema = z.array(notificationSchema);

export type NotificationSocketMessage = z.infer<
    typeof notificationSocketMessageSchema
>;
