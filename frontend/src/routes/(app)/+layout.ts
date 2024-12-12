import { authFetch } from '@/helpers/fetch';
import type { LayoutLoad } from './$types';
import type { User } from '@/types/user';
import { authToken } from '@/auth/token';
import { get } from 'svelte/store';
import { redirect } from '@sveltejs/kit';
import { getAllSensors } from '@/helpers/sensor';

const getUserData = async (fetchFN: typeof fetch) => {
    const res = await authFetch(`/api/v1/user/me`, {}, fetchFN);
    if (!res.ok) {
        throw new Error('cannot fetch current user');
    }
    const body = (await res.json()) as { user: User };

    return body.user;
};

export const load: LayoutLoad = async ({ fetch }) => {
    if (!get(authToken)) redirect(304, '/login');
    return {
        currentUser: await getUserData(fetch),
        sensors: await (async () => {
            const sensors = await getAllSensors(fetch);
            if (sensors.isError) return [];
            return sensors.data;
        })(),
    };
};
