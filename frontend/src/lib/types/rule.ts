import { z } from 'zod';




// --------------------- Non-Expandable ---------------------

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


// --------------------- Expandable ---------------------  

const RuleNot: z.ZodType<RuleNotType> = z.object({
    type: z.literal('not'),
    wrapped: z.lazy(() => ruleInternalSchema),
});

const RuleAnd: z.ZodType<RuleAndType> = z.object({
    type: z.literal('and'),
    children: z.lazy(() => ruleInternalSchema.array()),
});

const RuleOr: z.ZodType<RuleOrType> = z.object({
    type: z.literal('or'),
    children: z.lazy(() => ruleInternalSchema.array()),
});


// -------------------

type RuleNotType = {
    type: 'not';
    wrapped: RuleInternal;
};


type RuleAndType = {
    type: 'and';
    children: RuleInternal[];
};


type RuleOrType = {
    type: 'or';
    children: RuleInternal[];
};


export const ruleInternalSchema = z.union([
    RuleAnd,
    RuleOr,
    RuleNot,
    RuleGT,
    RuleLT,
]);

export type RuleInternal =
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
