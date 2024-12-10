import { authFetch } from '@/helpers/fetch';
import type { LayoutLoad } from './$types';
import type { User } from '@/types/user';
import { authToken } from '@/auth/token';
import { get } from 'svelte/store';
import { redirect } from '@sveltejs/kit';
import { getAllSensors } from '@/helpers/sensor';

const getUserData = async (fetchFN: typeof fetch) => {
    const res = await authFetch(`/api/v1/user`, {}, fetchFN);
    if (!res.ok) {
        return undefined;
    }
    const body = (await res.json()) as { user: User };

    return body.user;
};

export const load: LayoutLoad = async ({ fetch }) => {
    if (!get(authToken)) redirect(304, '/login');
    return {
        user: getUserData(fetch),
        sensors: await (async () => {
            const sensors = await getAllSensors(fetch);
            if (sensors.isError) return [];
            return sensors.data;
        })(),
    };
};
