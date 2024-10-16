import { authToken } from '@/auth/token';
import { get } from 'svelte/store';

type Url = RequestInfo | URL;
type Options = RequestInit | undefined;

export const authFetch = (
    url: Url,
    options: Options,
    fetchFun: (url: Url, options: Options) => Promise<Response> = fetch
) => {
    if (!options) options = {};
    if (!options.headers) options.headers = {};

    if (!('Authorization' in options.headers)) {
        const token = get(authToken);
        if (token) {
            options.headers = {
                ...options.headers,
                Authorization: `Bearer ${token.token}`,
            };
        }
    }
    return fetchFun(url, options);
};
