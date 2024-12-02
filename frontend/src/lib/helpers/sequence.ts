import type { FetchFn } from '@/types/misc';
import type { Result } from '@/types/result';
import { authFetch } from './fetch';
import { z } from 'zod';
import {
    sequenceDetailsSchema,
    sequenceSchema,
    type Sequence,
    type SequenceDetails,
} from '@/types/sequence';

export const getAllSequences = async (
    fetch: FetchFn
): Promise<Result<Sequence[], string>> => {
    const res = await authFetch(`/api/v1/sequence`, {}, fetch);
    const data = await res.json();
    if (!res.ok) {
        return {
            isError: true,
            error: data.error,
        };
    }

    const parsed = z.object({ data: sequenceSchema.array() }).safeParse(data);
    if (!parsed.success) {
        return {
            isError: true,
            error: 'Error while parsing the data! (getAllSequences)',
        };
    }

    return {
        isError: false,
        data: parsed.data.data,
    };
};

export const getSequenceDetails = async (
    id: string,
    fetch: FetchFn
): Promise<Result<SequenceDetails, string>> => {
    const res = await authFetch(`/api/v1/sequence/${id}`, {}, fetch);
    const data = await res.json();
    if (!res.ok) {
        return {
            isError: true,
            error: data.error,
        };
    }

    const parsed = z
        .object({ sequence: sequenceDetailsSchema })
        .safeParse(data);
    if (!parsed.success) {
        return {
            isError: true,
            error: 'Error while parsing the data! (getSequenceDetails)',
        };
    }

    return {
        isError: false,
        data: parsed.data.sequence,
    };
};

export const getAllDummySequences = async (): Promise<
    Result<Sequence[], string>
> => {
    return new Promise((resolve) => {
        resolve({
            isError: false,
            data: dummySequences,
        });
    });
};

export const getDummySequenceDetails = async (
    id: string
): Promise<Result<SequenceDetails, string>> => {
    return new Promise((resolve, reject) => {
        const sequence = dummySequenceDetails.find((s) => s.id === id);
        if (sequence)
            resolve({
                isError: false,
                data: sequence,
            });
        reject({
            isError: true,
            error: 'Not found',
        });
    });
};

// Dummy data for SequenceDetails
const dummySequenceDetails: SequenceDetails[] = [
    {
        id: '1e4a8d67-3b23-4fdb-95e5-a9d5ef0f4321',
        name: 'Morning Routine',
        description:
            'Automates lighting, heating, and coffee machine for the morning routine.',
        actions: [
            {
                target: 'device-light-livingroom',
                value: 100, // Set brightness to 100%
                msDelay: 500,
            },
            {
                target: 'device-heater-bedroom',
                value: 22, // Set bedroom temperature to 22°C
                msDelay: 2000,
            },
            {
                target: 'device-coffee-maker',
                value: 1, // Activate coffee machine
                msDelay: 3000,
            },
        ],
        created_at: new Date('2024-11-02T12:34:56Z'),
    },
    {
        id: '56fa3e94-9273-46b4-a9c3-8f87f0f9e123',
        name: 'Evening Security Check',
        description:
            'Activates security sensors and locks doors for nighttime security.',
        actions: [
            {
                target: 'sensor-door-lock-front',
                value: 1, // Lock front door
                msDelay: 500,
            },
            {
                target: 'sensor-window-lock-all',
                value: 1, // Lock all windows
                msDelay: 1500,
            },
            {
                target: 'sensor-motion-detector-outside',
                value: 1, // Activate motion detector
                msDelay: 2500,
            },
            {
                target: 'device-security-lights',
                value: 75, // Set outside lights to 75% brightness
                msDelay: 3500,
            },
        ],
        created_at: new Date('2024-11-02T12:34:56Z'),
    },
    {
        id: 'd9b437b1-5154-4be5-95d5-3349f7e6c123',
        name: 'Energy Saver Mode',
        description:
            'Reduces power usage by lowering heating and turning off unnecessary lights.',
        actions: [
            {
                target: 'device-heater-livingroom',
                value: 18, // Set temperature to 18°C
                msDelay: 1000,
            },
            {
                target: 'device-heater-bedroom',
                value: 18, // Set bedroom temperature to 18°C
                msDelay: 2000,
            },
            {
                target: 'device-light-kitchen',
                value: 0, // Turn off kitchen lights
                msDelay: 3000,
            },
            {
                target: 'device-light-hallway',
                value: 0, // Turn off hallway lights
                msDelay: 3500,
            },
        ],
        created_at: new Date('2024-11-02T12:34:56Z'),
    },
    {
        id: '2a1f9b5c-7c57-4c2b-8a8d-4b1f9e2d1234',
        name: 'Garden Watering Schedule',
        description:
            'Turns on water sprinklers and adjusts garden lights for evening watering.',
        actions: [
            {
                target: 'device-sprinkler-zone1',
                value: 1, // Activate sprinkler in zone 1
                msDelay: 500,
            },
            {
                target: 'device-sprinkler-zone2',
                value: 1, // Activate sprinkler in zone 2
                msDelay: 1500,
            },
            {
                target: 'device-light-garden',
                value: 50, // Set garden lights to 50% brightness
                msDelay: 2500,
            },
        ],
        created_at: new Date('2024-11-02T12:34:56Z'),
    },
    {
        id: '7f4d6a85-9d58-4c8a-8a6e-f3b9e7d5b123',
        name: 'Movie Time',
        description:
            'Adjusts lights and sound system for a cozy movie-watching experience.',
        actions: [
            {
                target: 'device-light-livingroom',
                value: 10, // Dim living room lights to 10%
                msDelay: 500,
            },
            {
                target: 'device-sound-system',
                value: 30, // Set sound system to volume level 30
                msDelay: 1000,
            },
            {
                target: 'device-tv',
                value: 1, // Turn on TV
                msDelay: 1500,
            },
        ],
        created_at: new Date('2024-11-02T12:34:56Z'),
    },
];
const dummySequences = dummySequenceDetails.map((s) => ({
    id: s.id,
    name: s.name,
    description: s.description,
}));
