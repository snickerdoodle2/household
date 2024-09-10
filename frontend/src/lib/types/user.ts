import { z } from "zod";

export const userSchema = z.object({
    id: z.string().uuid(),
    username: z.string(),
    name: z.string(),
    created_at: z.string().transform((d) => new Date(d)),
});

export type User = z.infer<typeof userSchema>;
