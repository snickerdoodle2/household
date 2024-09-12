import { z } from 'zod';
import { createZodStore } from '$lib/helpers/stores';
import { authFetch } from '@/helpers/fetch';
import { invalidateAll } from '$app/navigation';

const authTokenSchema = z.object({
    token: z.string(),
    expiry: z.string().transform((d) => new Date(d)),
});

const { set, subscribe } = createZodStore(authTokenSchema);

export const authToken = {
    subscribe,
    set: (v: unknown) => {
        const err = set(v);
        if (err) return err;
        localStorage.setItem('authToken', JSON.stringify(v));
        invalidateAll();
    },
    unset: async () => {
        set(undefined);
        localStorage.removeItem('authToken');
        await authFetch(`/api/v1/logout`, {
            method: 'POST',
        });
        invalidateAll();
    },
};
