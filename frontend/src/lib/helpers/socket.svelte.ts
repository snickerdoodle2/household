import { authToken } from "@/auth/token"
import { get } from "svelte/store"
import { z } from "zod"

const authSchema = z.object({
    type: z.literal('auth'),
    message: z.string()
})

const sensorDataSuccessSchema = z.object({
    status: z.literal('ok'),
    values: z.record(z.string().or(z.date()).transform(d => new Date(d)), z.number())
})

const sensorDataErrorSchema = z.object({
    status: z.literal('error'),
    message: z.string()
})

const subscribeSchema = z.object({
    type: z.literal('subscribe'),
    data: z.record(z.string(), z.discriminatedUnion('status', [sensorDataSuccessSchema, sensorDataErrorSchema]))
})

const messageSchema = z.discriminatedUnion('type', [authSchema, subscribeSchema])

export class SensorWebsocket {
    private toSubscribe: string[] = []
    private websocket: WebSocket
    data: Map<string, Map<Date, number>>

    constructor(toSubscribe: string[] | undefined = undefined) {
        toSubscribe ??= []
        this.toSubscribe = toSubscribe
        this.data = new Map()

        const token = get(authToken);
        if (!token) {
            throw new Error('Missing auth token!')
        }

        const url = new URL('/api/v1/sensor/measurements', window.location.href);
        url.protocol = url.protocol.replace('http', 'ws');
        this.websocket = new WebSocket(url, "inzynierka")

        // DEBUG
        this.websocket.addEventListener('message', e => {
            console.log(JSON.parse(e.data))
        })

        // AUTH
        this.websocket.addEventListener('message', e => {
            const { data: message, success } = messageSchema.safeParse(JSON.parse(e.data))
            if (!success) return;
            if (message.type === 'auth') {
                this.handleAuthMessage(message)
            }
            if (message.type === 'subscribe') {
                this.handleSubscribeMessage(message)
            }
        })


        this.websocket.addEventListener('open', () => {
            this.websocket.send(JSON.stringify({
                type: 'auth',
                data: token.token
            }))
        })

        this.websocket.addEventListener('error', (error) => {
            console.error(error)
        })
    }

    private handleAuthMessage(message: z.infer<typeof authSchema>) {
        if (message.message !== 'ok') {
            // TODO: do smth with it
            throw new Error('not ok');
        }

        if (this.toSubscribe.length > 0) {
            this.websocket.send(JSON.stringify({
                type: "subscribe",
                data: this.toSubscribe
            }))
        }
    }

    private handleSubscribeMessage(message: z.infer<typeof subscribeSchema>) {
        console.log(message.data)
    }
}
