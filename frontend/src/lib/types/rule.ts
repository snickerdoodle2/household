import { z } from "zod";

const RuleGT = z.object({
    type: z.literal("gt"),
    sensor_id: z.string().uuid(),
    value: z.number(),
});

const RuleLT = z.object({
    type: z.literal("lt"),
    sensor_id: z.string().uuid(),
    value: z.number(),
});

type RuleAndType = {
    type: "and";
    children: RuleInternalType[];
};

const RuleAnd: z.ZodType<RuleAndType> = z.object({
    type: z.literal("and"),
    children: z.lazy(() => RuleInternal.array()),
});

type RuleOrType = {
    type: "or";
    children: RuleInternalType[];
};

const RuleOr: z.ZodType<RuleOrType> = z.object({
    type: z.literal("or"),
    children: z.lazy(() => RuleInternal.array()),
});

const RuleInternal = z.union([RuleAnd, RuleOr, RuleGT, RuleLT]);

type RuleInternalType =
    | RuleAndType
    | RuleOrType
    | z.infer<typeof RuleGT>
    | z.infer<typeof RuleLT>;

export const Rule = z.object({
    id: z.string().uuid(),
    description: z.string().max(256),
    on_valid: z.object({
        to: z.string().uuid(),
        payload: z.object({}).passthrough(),
    }),
    internal: RuleInternal,
});

export type RuleType = z.infer<typeof Rule>;
