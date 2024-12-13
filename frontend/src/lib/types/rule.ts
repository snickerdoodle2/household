import { z } from 'zod';

const RuleGT = z.object({
    type: z.literal('gt'),
    sensor_id: z.string().uuid(),
    value: z.number(),
});

export type RuleGtType = z.infer<typeof RuleGT>;

const RuleLT = z.object({
    type: z.literal('lt'),
    sensor_id: z.string().uuid(),
    value: z.number(),
});

export type RuleLtType = z.infer<typeof RuleLT>;

const RulePerc = z.object({
    type: z.literal('perc'),
    sensor_id: z.string().uuid(),
    duration: z
        .string()
        .regex(/^-?(?:\d+(?:\.\d+)?(?:h|m|s|(?:ms)|(?:µs)|(?:us)|(?:ns)))+$/),
    perc: z.number().int().min(0).max(100),
});

export type RulePercType = z.infer<typeof RulePerc>;

const RuleTime = z.object({
    type: z.literal('time'),
    variant: z.literal('before').or(z.literal('after')),
    hour: z.number().min(0).max(23),
    minute: z.number().min(0).max(59),
});

export type RuleTimeType = z.infer<typeof RulePerc>;

const RuleNot: z.ZodType<RuleNotType> = z.object({
    type: z.literal('not'),
    wrapped: z.lazy(() => ruleInternalSchema),
});

export type RuleNotType = {
    type: 'not';
    wrapped: RuleInternal;
};

const RuleAnd: z.ZodType<RuleAndType> = z.object({
    type: z.literal('and'),
    children: z.lazy(() => ruleInternalSchema.array()),
});

export type RuleAndType = {
    type: 'and';
    children: RuleInternal[];
};

const RuleOr: z.ZodType<RuleOrType> = z.object({
    type: z.literal('or'),
    children: z.lazy(() => ruleInternalSchema.array()),
});

export type RuleOrType = {
    type: 'or';
    children: RuleInternal[];
};

export const ruleInternalSchema = z.union([
    RuleAnd,
    RuleOr,
    RuleNot,
    RuleGT,
    RuleLT,
    RulePerc,
    RuleTime,
]);

export type RuleInternal = z.infer<typeof ruleInternalSchema>;

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

const targetTypeEnum = z.enum(['sensor', 'sequence']);

const internalRuleSchema = z.object({
    on_valid: z.object({
        target_type: targetTypeEnum,
        target_id: z.string().uuid(),
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
