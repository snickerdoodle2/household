import { authToken } from '@/auth/token';
import { socketMessageSchema, type SocketMessage } from '@/types/socketMessage';
import { get, writable } from 'svelte/store';

export const socketStore = (id: string) => {
    const token = get(authToken);
    if (!token) {
        throw new Error('auth token is required');
    }

    const { set, subscribe } = writable<SocketMessage | null>();

    const url = new URL(`/api/v1/sensor/${id}/value`, window.location.href);
    url.protocol = url.protocol.replace('http', 'ws');
    url.searchParams.set('token', token.token);

    const socket = new WebSocket(url.toString());

    socket.addEventListener('message', (message) => {
        const { data, success } = socketMessageSchema.safeParse(
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
