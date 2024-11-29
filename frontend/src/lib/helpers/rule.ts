import type { FetchFn } from '@/types/misc';
import type { Result } from '@/types/result';
import {
    ruleDetailsSchema,
    ruleSchema,
    type Rule,
    type RuleDetails,
} from '@/types/rule';
import { authFetch } from './fetch';
import { z } from 'zod';

export const RULE_URL = '/api/v1/rule';

export const getAllRules = async (
    fetch: FetchFn
): Promise<Result<Rule[], string>> => {
    const res = await authFetch(RULE_URL, {}, fetch);
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

export const getRuleDetails = async (
    id: string,
    fetch: FetchFn
): Promise<Result<RuleDetails, string>> => {
    const res = await authFetch(`${RULE_URL}/${id}`, {}, fetch);
    const data = await res.json();
    if (!res.ok) {
        return {
            isError: true,
            error: data.error,
        };
    }

    const parsed = z.object({ rule: ruleDetailsSchema }).safeParse(data);
    if (!parsed.success) {
        return {
            isError: true,
            error: 'Error while parsing the data! (getRuleDetails)',
        };
    }

    return {
        isError: false,
        data: parsed.data.rule,
    };
};
