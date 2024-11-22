import { getAllUsers } from '@/helpers/user';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
    return {
        users: (async () => {
            const res = await getAllUsers(fetch);
            if (res.isError) {
                throw res.error;
            }
            console.log(res.data);
            return res.data;
        })(),
    };
};
