import { z } from "zod";
import { createZodStore } from "./helpers/stores";

const authTokenSchema = z.object({
    token: z.string(),
    expiry: z.string().transform((d) => new Date(d)),
});

export const authToken = createZodStore(authTokenSchema);

const userSchema = z.object({
    id: z.string().uuid(),
    username: z.string(),
    name: z.string(),
    createdAt: z.string().transform((d) => new Date(d)),
});

export const user = createZodStore(userSchema);
