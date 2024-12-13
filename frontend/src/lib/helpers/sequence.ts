import type { FetchFn } from '@/types/misc';
import type { Result } from '@/types/result';
import { authFetch } from './fetch';
import { z } from 'zod';
import {
    sequenceDetailsSchema,
    sequenceSchema,
    type Sequence,
    type SequenceDetails,
} from '@/types/sequence';

export const getAllSequences = async (
    fetch: FetchFn
): Promise<Result<Sequence[], string>> => {
    const res = await authFetch(`/api/v1/sequence`, {}, fetch);
    const data = await res.json();
    if (!res.ok) {
        return {
            isError: true,
            error: data.error,
        };
    }

    const parsed = z.object({ data: sequenceSchema.array() }).safeParse(data);
    if (!parsed.success) {
        return {
            isError: true,
            error: 'Error while parsing the data! (getAllSequences)',
        };
    }

    return {
        isError: false,
        data: parsed.data.data,
    };
};

export const getSequenceDetails = async (
    id: string,
    fetch: FetchFn
): Promise<Result<SequenceDetails, string>> => {
    const res = await authFetch(`/api/v1/sequence/${id}`, {}, fetch);
    const data = await res.json();
    if (!res.ok) {
        return {
            isError: true,
            error: data.error,
        };
    }

    const parsed = z
        .object({ sequence: sequenceDetailsSchema })
        .safeParse(data);
    if (!parsed.success) {
        return {
            isError: true,
            error: 'Error while parsing the data! (getSequenceDetails)',
        };
    }

    return {
        isError: false,
        data: parsed.data.sequence,
    };
};
