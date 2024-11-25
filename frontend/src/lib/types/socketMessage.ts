import { z } from 'zod';

export const socketMessageSchema = z.object({
    values: z.number().array(),
    status: z.enum(['ONLINE', 'OFFLINE']),
});

export type SocketMessage = z.infer<typeof socketMessageSchema>;
