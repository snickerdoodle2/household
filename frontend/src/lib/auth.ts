import { writable } from "svelte/store";
import { z } from "zod";

export const AuthTokenSchema = z.object({
    token: z.string(),
    expiry: z.string().transform((d) => new Date(d)),
});
export type AuthToken = z.infer<typeof AuthTokenSchema>;

export const authToken = writable<AuthToken | undefined>(undefined);

export const UserSchema = z.object({
    id: z.string().uuid(),
    username: z.string(),
    name: z.string(),
    createdAt: z.string().transform((d) => new Date(d)),
});
export type User = z.infer<typeof UserSchema>;

export const user = writable<AuthToken | undefined>(undefined);
