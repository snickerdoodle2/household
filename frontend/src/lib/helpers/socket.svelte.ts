import { authToken } from '@/auth/token';
import { SvelteMap } from 'svelte/reactivity';
import { get } from 'svelte/store';
import { z } from 'zod';

const durationSchema = z
    .string()
    .regex(/^-?(?:\d+(?:\.\d+)?(?:h|m|s|(?:ms)|(?:µs)|(?:us)|(?:ns)))+$/);

const authSchema = z.object({
    type: z.literal('auth'),
    message: z.string(),
});

const sensorDataSuccessSchema = z.object({
    status: z.literal('ok'),
    values: z
        .record(z.string().datetime({ offset: true }), z.number())
        .transform((e) => {
            return Object.entries(e).reduce((acc, [key, value]) => {
                acc.set(new Date(key), value);
                return acc;
            }, new SvelteMap<Date, number>());
        }),
});

const sensorDataErrorSchema = z.object({
    status: z.literal('error'),
    message: z.string(),
});

const measurementResponseSchema = z.object({
    type: z.literal('measurement_req'),
    id: z.string(),
    values: z
        .record(z.string().datetime({ offset: true }), z.number())
        .transform((e) => {
            return Object.entries(e).reduce((acc, [key, value]) => {
                acc.set(new Date(key), value);
                return acc;
            }, new SvelteMap<Date, number>());
        }),
});

const subscribeSchema = z.object({
    type: z.literal('subscribe'),
    data: z.record(
        z.string(),
        z.discriminatedUnion('status', [
            sensorDataSuccessSchema,
            sensorDataErrorSchema,
        ])
    ),
});

const measurementSchema = z.object({
    type: z.literal('measurment'),
    sensor_id: z.string().uuid(),
    time: z
        .string()
        .datetime({ offset: true })
        .or(z.date())
        .transform((d) => new Date(d)),
    value: z.number(),
});

// HANDLE UNSUBSCRIBE MESSAGE
const messageSchema = z.discriminatedUnion('type', [
    authSchema,
    subscribeSchema,
    measurementSchema,
    measurementResponseSchema,
]);

export class AppWebsocket {
    private websocket!: WebSocket;
    private subscriptionCount!: Map<string, number>;
    ready = $state(false);
    data: SvelteMap<string, SvelteMap<Date, number>> = $state(new SvelteMap());
    private static _instance: AppWebsocket | null = null;

    constructor() {
        if (AppWebsocket._instance) {
            return AppWebsocket._instance;
        }

        AppWebsocket._instance = this;

        this.subscriptionCount = new Map();

        const token = get(authToken);
        if (!token) {
            throw new Error('Missing auth token!');
        }

        const url = new URL(
            '/api/v1/sensor/measurements',
            window.location.href
        );
        url.protocol = url.protocol.replace('http', 'ws');
        this.websocket = new WebSocket(url, 'inzynierka');

        // DEBUG
        this.websocket.addEventListener('message', (e) => {
            console.log(JSON.parse(e.data));
        });

        // AUTH
        this.websocket.addEventListener('message', (e) => {
            const {
                data: message,
                success,
                error,
            } = messageSchema.safeParse(JSON.parse(e.data));
            if (!success) {
                console.error(error.issues);
                return;
            }
            if (message.type === 'auth') {
                this.handleAuthMessage(message);
            }
            if (message.type === 'subscribe') {
                this.handleSubscribeMessage(message);
            }
            if (message.type === 'measurment') {
                this.handleMeasurementMessage(message);
            }
            if (message.type === 'measurement_req') {
                this.handleMeasurementResponse(message);
            }
        });

        this.websocket.addEventListener('open', () => {
            this.websocket.send(
                JSON.stringify({
                    type: 'auth',
                    data: token.token,
                })
            );
            this.ready = true;
        });

        this.websocket.addEventListener('error', (error) => {
            console.error(error);
        });
    }

    private handleAuthMessage(message: z.infer<typeof authSchema>) {
        if (message.message !== 'ok') {
            // TODO: do smth with it
            throw new Error('not ok');
        }
    }

    private handleSubscribeMessage(message: z.infer<typeof subscribeSchema>) {
        for (const [key, value] of Object.entries(message.data)) {
            if (value.status === 'error') continue;
            this.data.set(key, value.values);
        }
    }

    private handleMeasurementResponse(
        message: z.infer<typeof measurementResponseSchema>
    ) {
        this.data.set(message.id, message.values);
    }

    subscribe(sensorID: string) {
        let cur = this.subscriptionCount.get(sensorID) ?? 0;
        cur += 1;
        this.subscriptionCount.set(sensorID, cur);

        // if not subscribed already
        if (cur <= 1) {
            this.websocket.send(
                JSON.stringify({
                    type: 'subscribe',
                    data: [sensorID],
                })
            );
        }
    }

    unsubscribe(sensorID: string) {
        let cur = this.subscriptionCount.get(sensorID) ?? 0;

        // return if no one is subscribed
        if (cur <= 0) {
            return;
        }
        cur -= 1;
        this.subscriptionCount.set(sensorID, cur);

        // send unsubscribe msg if no one is subscribed
        if (cur === 0) {
            this.websocket.send(
                JSON.stringify({
                    type: 'unsubscribe',
                    data: sensorID,
                })
            );
            this.data.delete(sensorID);
        }
    }

    public requestSince(sensorID: string, delta: string) {
        const { success, data: duration } = durationSchema.safeParse(delta);
        if (!success) {
            throw new Error('Invalid delta format');
        }

        this.websocket.send(
            JSON.stringify({
                type: 'measurement_req',
                data: {
                    id: sensorID,
                    duration,
                },
            })
        );
    }

    close() {
        this.websocket.close();
    }

    private handleMeasurementMessage(
        message: z.infer<typeof measurementSchema>
    ) {
        this.data.get(message.sensor_id)?.set(message.time, message.value);
    }
}
