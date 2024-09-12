import { authFetch } from '@/helpers/fetch';
import type { LayoutLoad } from './$types';
import type { User } from '@/types/user';
import { authToken } from '@/auth/token';
import { get } from 'svelte/store';

const getUserData = async (fetchFN: typeof fetch) => {
    const token = get(authToken);
    if (!token) return undefined;
    const res = await authFetch(`/api/v1/user`, {}, fetchFN);
    if (!res.ok) {
        return undefined;
    }
    const body = (await res.json()) as { user: User };

    return body.user;
};

export const load: LayoutLoad = async ({ fetch }) => {
    return {
        user: await getUserData(fetch),
    };
};
