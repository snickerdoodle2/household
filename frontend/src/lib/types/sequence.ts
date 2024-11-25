import { z } from 'zod';

export const sequenceActionSchema = z.object({
    target: z.string().uuid(),
    value: z.number().gte(0).finite(),
    msDelay: z.number().gte(0).finite(), // delay from start of the sequence given in ms
});

export const sequenceSchema = z.object({
    id: z.string().uuid(),
    name: z.string().min(1).max(32),
    description: z.string().max(256),
});

export const sequenceDetailsSchema = sequenceSchema.merge(
    z.object({
        actions: z.array(sequenceActionSchema),
        created_at: z
            .string()
            .or(z.date())
            .transform((d) => new Date(d)),
    })
);

export const newSequenceSchema = sequenceDetailsSchema.omit({
    id: true,
    created_at: true,
});

export type Sequence = z.infer<typeof sequenceSchema>;
export type SequenceAction = z.infer<typeof sequenceActionSchema>;
export type SequenceDetails = z.infer<typeof sequenceDetailsSchema>;
export type NewSequence = z.infer<typeof newSequenceSchema>;
