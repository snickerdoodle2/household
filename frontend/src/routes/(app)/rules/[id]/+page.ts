import { getRuleDetails } from '@/helpers/rule';
import type { PageLoad } from './$types';
import { getAllSequences } from '@/helpers/sequence';

export const load: PageLoad = async ({ params, fetch }) => {
    return {
        rule: (async () => {
            const res = await getRuleDetails(params.id, fetch);
            if (res.isError) {
                throw res.error;
            }
            return res.data;
        })(),
        sequences: (async () => {
            const sequences = await getAllSequences(fetch);
            if (sequences.isError) return [];
            return sequences.data;
        })(),
    };
};
