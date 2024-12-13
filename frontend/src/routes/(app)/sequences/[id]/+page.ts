import type { PageLoad } from './$types';
import { getSequenceDetails } from '@/helpers/sequence';

export const load: PageLoad = async ({ params, fetch }) => {
    return {
        sequence: (async () => {
            const res = await getSequenceDetails(params.id, fetch);
            if (res.isError) {
                throw res.error;
            }
            return res.data;
        })(),
    };
};
