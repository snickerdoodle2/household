import { get } from 'svelte/store';
import type { PageLoad } from './$types';
import { authToken } from '@/auth/token';
import { redirect } from '@sveltejs/kit';

export const load: PageLoad = () => {
    if (get(authToken)) redirect(301, '/');
};
