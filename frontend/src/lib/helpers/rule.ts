import type { FetchFn } from '@/types/misc';
import type { Result } from '@/types/result';
import { ruleSchema, type Rule } from '@/types/rule';
import { authFetch } from './fetch';
import { z } from 'zod';

export const getAllRules = async (
    fetch: FetchFn
): Promise<Result<Rule[], string>> => {
    const res = await authFetch(`/api/v1/rule`, {}, fetch);
    const data = await res.json();
    if (!res.ok) {
        return {
            isError: true,
            error: data.error,
        };
    }

    const parsed = z.object({ data: ruleSchema.array() }).safeParse(data);
    if (!parsed.success) {
        return {
            isError: true,
            error: 'Error while parsing the data! (getAllRules)',
        };
    }

    return {
        isError: false,
        data: parsed.data.data,
    };
};
