import { invalidateAll } from '$app/navigation';
import { authFetch } from '@/helpers/fetch';
import type { Login } from '@/types/login';
import { persisted } from 'svelte-persisted-store';
import { z } from 'zod';

const authTokenSchema = z.object({
    token: z.string(),
    expiry: z.string().transform((d) => new Date(d)),
});

type AuthToken = z.infer<typeof authTokenSchema>;

export const authToken = (() => {
    const { set, subscribe } = persisted<AuthToken | null>('authToken', null);

    return {
        subscribe,
        login: async (data: Login) => {
            const res = await fetch('/api/v1/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
            });

            if (!res.ok) {
                return await res.json();
            }

            const token_ = (await res.json())['auth_token'];

            const { data: token, success } = authTokenSchema.safeParse(token_);
            if (!success) return {};
            set(token);
        },
        logout: async () => {
            try {
                await authFetch('/api/v1/logout', {
                    method: 'POST',
                });
            } finally {
                set(null);
                invalidateAll();
            }
        },
    };
})();
