import { authToken } from '@/auth/token';
import {
    sensorSocketMessageSchema,
    type SensorSocketMessage,
} from '@/types/socketMessage';
import { get, writable } from 'svelte/store';

export const socketStore = (id: string) => {
    const token = get(authToken);
    if (!token) {
        throw new Error('auth token is required');
    }

    const { set, subscribe } = writable<SensorSocketMessage | null>();

    const url = new URL(`/api/v1/sensor/${id}/value`, window.location.href);
    url.protocol = url.protocol.replace('http', 'ws');
    url.searchParams.set('token', token.token);

    const socket = new WebSocket(url.toString());

    socket.addEventListener('message', (message) => {
        const { data, success } = sensorSocketMessageSchema.safeParse(
            JSON.parse(message.data)
        );
        if (!success) return;

        set(data);
    });

    return {
        subscribe,
        close: () => {
            socket.close();
        },
    };
};

export type SocketStore = ReturnType<typeof socketStore>;
