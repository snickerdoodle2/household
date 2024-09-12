import { authFetch } from '@/helpers/fetch';
import { persisted } from 'svelte-persisted-store';
import { z } from 'zod';

const authTokenSchema = z.object({
    token: z.string(),
    expiry: z.string().transform((d) => new Date(d)),
});

type AuthToken = z.infer<typeof authTokenSchema>;

export const authToken = (() => {
    const { set, subscribe, reset } = persisted<AuthToken | null>(
        'authToken',
        null
    );

    return {
        subscribe,
        set: (value: unknown) => {
            const { data, success, error } = authTokenSchema.safeParse(value);
            if (!success) {
                reset();
                return error.issues;
            }
            set(data);
        },
        logout: async () => {
            try {
                await authFetch('/api/v1/logout', {
                    method: 'POST',
                });
            } finally {
                reset();
            }
        },
    };
})();
