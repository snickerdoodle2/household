import { getRuleDetails } from '@/helpers/rule';
import type { PageLoad } from './$types';
import { getAllSensors } from '@/helpers/sensor';

export const load: PageLoad = async ({ params, fetch }) => {
    return {
        rule: (async () => {
            const res = await getRuleDetails(params.id, fetch);
            if (res.isError) {
                throw res.error;
            }
            return res.data;
        })(),
        sensors: (async () => {
            const res = await getAllSensors(fetch);
            if (res.isError) {
                throw res.error;
            }
            return res.data;
        })(),
    };
};
