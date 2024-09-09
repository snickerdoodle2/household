import { z } from "zod";
import { createZodStore } from "$lib/helpers/stores";
import { browser } from "$app/environment";

const authTokenSchema = z.object({
    token: z.string(),
    expiry: z.string().transform((d) => new Date(d)),
});

const token = browser
    ? (JSON.parse(localStorage.getItem("authToken") ?? "") as z.infer<
          typeof authTokenSchema
      >)
    : undefined;
const { set, subscribe } = createZodStore(authTokenSchema, token);

export const authToken = {
    subscribe,
    set: (v: any) => {
        const err = set(v);
        if (err) return err;
        localStorage.setItem("authToken", JSON.stringify(v));
    },
};
