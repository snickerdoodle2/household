import { getAllSensors } from '@/helpers/sensor';
import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { get } from 'svelte/store';
import { authToken } from '@/auth/token';

export const load: PageLoad = async ({ fetch }) => {
    if (!get(authToken)) return undefined;
    const sensors = await getAllSensors(fetch);

    if (sensors.isError) {
        error(500, sensors.error);
    }

    return {
        sensors: sensors.data.map((e) => {
            return { label: e.name, value: e.id };
        }),
    };
};
