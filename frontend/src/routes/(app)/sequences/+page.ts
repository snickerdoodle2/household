import type { PageLoad } from './$types';
import { getAllSequences } from '@/helpers/sequence';

export const load: PageLoad = async () => {
    return {
        sequences: (async () => {
            const sequences = await getAllSequences(fetch);
            if (sequences.isError) return [];
            return sequences.data;
        })(),
    };
};
