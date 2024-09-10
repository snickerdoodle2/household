import { z } from "zod";
import { createZodStore } from "$lib/helpers/stores";
import { browser } from "$app/environment";
import { authFetch } from "@/helpers/fetch";
import { SERVER_URL } from "@/const";
import { invalidateAll } from "$app/navigation";

const authTokenSchema = z.object({
    token: z.string(),
    expiry: z.string().transform((d) => new Date(d)),
});

let token = undefined;
if (browser) {
    const tokenText = localStorage.getItem("authToken");
    if (tokenText) {
        token = JSON.parse(tokenText) as z.infer<typeof authTokenSchema>;
    }
}

const { set, subscribe } = createZodStore(authTokenSchema, token);

export const authToken = {
    subscribe,
    set: (v: any) => {
        const err = set(v);
        if (err) return err;
        localStorage.setItem("authToken", JSON.stringify(v));
        invalidateAll();
    },
    unset: async () => {
        set(undefined);
        localStorage.removeItem("authToken");
        await authFetch(`${SERVER_URL}/api/v1/logout`, {
            method: "POST",
        });
        invalidateAll();
    },
};
