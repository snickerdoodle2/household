import type { PageLoad } from './$types';
import { getDummySequenceDetails } from '@/helpers/sequence';

export const load: PageLoad = async ({ params }) => {
    return {
        sequence: (async () => {
            const res = await getDummySequenceDetails(params.id);
            if (res.isError) {
                throw res.error;
            }
            return res.data;
        })(),
    };
};
