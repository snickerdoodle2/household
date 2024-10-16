import { getAllRules } from '@/helpers/rule';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
    return {
        rules: (async () => {
            const rules = await getAllRules(fetch);
            if (rules.isError) return [];
            return rules.data;
        })(),
    };
};
