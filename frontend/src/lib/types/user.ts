import { z } from 'zod';

export const userSchema = z.object({
    id: z.string().uuid().optional(),
    username: z.string(),
    name: z.string(),
    created_at: z
        .string()
        .transform((d) => new Date(d))
        .optional(),
});

export const newUserSchema = userSchema
    .omit({ id: true, created_at: true })
    .merge(
        z.object({
            password: z.string().min(8).max(32),
            confirmPassword: z.string().min(8).max(32),
        })
    )
    .refine((e) => e.password === e.confirmPassword, {
        message: 'Passwords must match!',
        path: ['confirmPassword'],
    });

export type User = z.infer<typeof userSchema>;
export type NewUser = z.infer<typeof newUserSchema>;
