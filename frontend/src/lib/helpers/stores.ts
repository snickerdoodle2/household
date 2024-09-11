import { writable } from 'svelte/store';
import type { z } from 'zod';

export function createZodStore<T extends z.ZodTypeAny>(
    schema: T,
    defaultValue: z.infer<T> | undefined = undefined
) {
    type Type = z.infer<T>;
    const { subscribe, set } = writable<Type | undefined>(defaultValue);

    return {
        subscribe,
        set: (v: any) => {
            if (v === undefined) {
                set(undefined);
                return undefined;
            }
            const { data, success, error } = schema.safeParse(v);
            if (!success) {
                set(undefined);
                return error;
            }
            set(data as Type);
            return undefined;
        },
    };
}
