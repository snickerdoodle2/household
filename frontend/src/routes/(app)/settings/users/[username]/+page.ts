import { getUser } from '@/helpers/user';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
    return {
        user: await getUser(params.username, fetch),
    };
};
