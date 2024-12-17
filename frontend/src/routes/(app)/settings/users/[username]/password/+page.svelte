<script lang="ts">
    import { Button } from '@/components/ui/button';
    import { Input } from '$lib/components/ui/input';
    import { Label } from '$lib/components/ui/label';
    import { z } from 'zod';
    import { authFetch } from '@/helpers/fetch';
    import type { PageData } from './$types';
    import { goto } from '$app/navigation';

    type Props = {
        data: PageData;
    };

    const { data }: Props = $props();

    let password = $state('');
    let confirmPassword = $state('');
    let errors = $state('');
    const passwordSchema = z
        .object({
            password: z.string().min(8).max(32),
            confirmPassword: z.string().min(8).max(32),
        })
        .refine((e) => e.password === e.confirmPassword, {
            path: ['confirmPassword'],
            message: 'Password must match',
        });

    const onsubmit = async () => {
        const {
            data: parseData,
            success,
            error,
        } = passwordSchema.safeParse({
            password,
            confirmPassword,
        });

        if (!success) {
            errors = JSON.stringify(error.issues);
            return;
        }

        const res = await authFetch(`/api/v1/user/${data.username}`, {
            method: 'PUT',
            body: JSON.stringify({
                password: parseData.password,
            }),
        });

        if (res.ok) {
            goto('/settings/users');
        }
    };
</script>

<div class="mx-auto max-w-[256px]">
    <form class="flex flex-col gap-4" {onsubmit}>
        <fieldset>
            <Label for="password">Password</Label>
            <Input type="password" bind:value={password} name="password" />
        </fieldset>
        <fieldset>
            <Label for="confirmPassword">Confirm password</Label>
            <Input
                type="password"
                bind:value={confirmPassword}
                name="confirmPassword"
            />
        </fieldset>
        {#if errors.length}
            <!-- TODO: ładniejszy error handling -->
            <code class="text-red-500">{errors}</code>
        {/if}
        <Button type="submit">Change Password</Button>
    </form>
</div>
