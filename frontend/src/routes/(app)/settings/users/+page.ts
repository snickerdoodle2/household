import { getAllUsers } from '@/helpers/user';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
    return {
        users: await getAllUsers(fetch),
    };
};
