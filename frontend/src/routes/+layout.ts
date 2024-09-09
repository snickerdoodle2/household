import { authFetch } from "@/helpers/fetch";
import type { LayoutLoad } from "./$types";
import { SERVER_URL } from "@/const";
import type { User } from "@/types/user";

export const ssr = false;

const getUserData = async (fetchFN: typeof fetch) => {
    const res = await authFetch(`${SERVER_URL}/api/v1/user`, {}, fetchFN);
    if (!res.ok) {
        return undefined;
    }
    const body = (await res.json()) as { user: User };

    return body.user;
};

export const load: LayoutLoad = async ({ fetch }) => {
    return {
        user: await getUserData(fetch),
    };
};
