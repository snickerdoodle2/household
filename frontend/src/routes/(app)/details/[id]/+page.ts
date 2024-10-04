import type { PageLoad } from './$types';
import { getSensorDetails } from '@/helpers/sensor';

export const load: PageLoad = async ({ params, fetch }) => {
    return {
        sensor: (async () => {
            const res = await getSensorDetails(params.id, fetch);
            if (res.isError) {
                throw res.error;
            }
            console.log(res.data);
            return res.data;
        })(),
    };
};