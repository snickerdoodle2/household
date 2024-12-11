import type { FetchFn } from '@/types/misc';
import { userSchema, type User } from '@/types/user';
import { authFetch } from './fetch';
import { z } from 'zod';

export const getAllUsers = async (fetch: FetchFn): Promise<User[]> => {
    const res = await authFetch('/api/v1/user', {}, fetch);
    const body = await res.json();
    if (!res.ok) {
        console.error(body);
        return [];
    }

    const { success, data, error } = z
        .object({ data: userSchema.array() })
        .safeParse(body);

    if (!success) {
        console.error(error.issues);
        return [];
    }

    return data.data;
};

export const getUser = async (
    username: string,
    fetch: FetchFn
): Promise<User | null> => {
    const res = await authFetch(`/api/v1/user/${username}`, {}, fetch);
    const body = await res.json();
    if (!res.ok) {
        console.error(body);
        return null;
    }

    const { success, data, error } = z
        .object({ user: userSchema })
        .safeParse(body);

    if (!success) {
        console.error(error.issues);
        return null;
    }

    return data.user;
};
