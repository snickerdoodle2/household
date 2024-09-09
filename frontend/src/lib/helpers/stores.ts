import { writable } from "svelte/store";
import type { z } from "zod";

export function createZodStore<T extends z.ZodTypeAny>(schema: T) {
    type Type = z.infer<T>;
    const { subscribe, set } = writable<Type | undefined>(undefined);

    return {
        subscribe,
        set: (v: any) => {
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
