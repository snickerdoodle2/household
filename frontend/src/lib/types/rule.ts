import { z } from 'zod';

const RuleGT = z.object({
    type: z.literal('gt'),
    sensor_id: z.string().uuid(),
    value: z.number(),
});

const RuleLT = z.object({
    type: z.literal('lt'),
    sensor_id: z.string().uuid(),
    value: z.number(),
});

type RuleNotType = {
    type: 'not';
    wrapped: RuleInternal;
};

const RuleNot: z.ZodType<RuleNotType> = z.object({
    type: z.literal('not'),
    wrapped: z.lazy(() => ruleInternalSchema),
});

type RuleAndType = {
    type: 'and';
    children: RuleInternal[];
};

const RuleAnd: z.ZodType<RuleAndType> = z.object({
    type: z.literal('and'),
    children: z.lazy(() => ruleInternalSchema.array()),
});

type RuleOrType = {
    type: 'or';
    children: RuleInternal[];
};

const RuleOr: z.ZodType<RuleOrType> = z.object({
    type: z.literal('or'),
    children: z.lazy(() => ruleInternalSchema.array()),
});

export const ruleInternalSchema = z.union([
    RuleAnd,
    RuleOr,
    RuleNot,
    RuleGT,
    RuleLT,
]);

type RuleInternal =
    | RuleAndType
    | RuleOrType
    | RuleNotType
    | z.infer<typeof RuleGT>
    | z.infer<typeof RuleLT>;

const ruleNameDescSchema = z.object({
    name: z.string().min(1).max(32),
    description: z.string().max(256),
});

export const ruleSchema = ruleNameDescSchema.merge(
    z.object({
        id: z.string().uuid(),
    })
);

export type Rule = z.infer<typeof ruleSchema>;

const internalRuleSchema = z.object({
    on_valid: z.object({
        to: z.string().uuid(),
        payload: z.object({}).passthrough(),
    }),
    internal: ruleInternalSchema,
});

export const ruleDetailsSchema = ruleSchema.merge(internalRuleSchema).merge(
    z.object({
        created_at: z
            .string()
            .or(z.date())
            .transform((d) => new Date(d)),
    })
);

export type RuleDetails = z.infer<typeof ruleDetailsSchema>;

export const newRuleSchema = internalRuleSchema.merge(ruleNameDescSchema);
export type NewRule = z.infer<typeof newRuleSchema>;
