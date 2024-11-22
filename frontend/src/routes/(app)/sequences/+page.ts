import type { PageLoad } from './$types';
import { getAllDummySequences } from '@/helpers/sequence';

export const load: PageLoad = async () => {
    return {
        sequences: (async () => {
            // const sequences = await getAllSequences(fetch);
            const sequences = await getAllDummySequences();
            if (sequences.isError) return [];
            return sequences.data;
        })(),
    };
};
