<script lang="ts">
    import { userSchema, type User } from '@/types/user';
    import { Label } from '$lib/components/ui/label';
    import * as Select from '$lib/components/ui/select';
    import NewSensorInput from '$lib/components/FormInput.svelte';
    import Button from '$lib/components/ui/button/button.svelte';
    import { authFetch } from '@/helpers/fetch';
    import { goto } from '$app/navigation';
    import type { Selected } from 'bits-ui';

    type Props = {
        user: User;
        currentUser: User;
    };

    const props: Props = $props();
    const ROLES = [
        {
            label: 'Admin',
            value: 'admin',
        },
        {
            label: 'User',
            value: 'user',
        },
    ] as const;

    let selectedRole = $state<Selected<'admin' | 'user'>>(
        ROLES.find((e) => e.value === props.user.role) ?? ROLES[0]
    );

    let user = $state(props.user);

    const cancelEditing = () => {
        user = props.user;
        editing = false;
    };

    let editing = $state(false);

    let globalError: string | undefined = $state();
    let fieldErrors: Record<string, string> = $state({});

    async function handleEdit() {
        const { data, success, error } = userSchema.safeParse({
            id: user.id,
            role: user.role,
            username: user.username,
            name: user.name,
        });

        if (!success) {
            console.log(error.issues);
            fieldErrors = {};
            for (const issue of error.issues) {
                fieldErrors[issue.path[0]] = issue.message;
            }
            return;
        }

        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        const { username, id, ...userData } = data;

        const res = await authFetch(`/api/v1/user/${username}`, {
            method: 'PUT',
            body: JSON.stringify(userData),
        });

        const resJson = await res.json();

        if (typeof resJson.error === 'string') {
            globalError = resJson.error;
        } else {
            fieldErrors = {};
        }

        if (res.ok) {
            editing = false;
        }
    }

    const handleDelete = async () => {
        const res = await authFetch(`/api/v1/user/${props.user.username}`, {
            method: 'DELETE',
        });

        if (res.ok) {
            goto('/settings/users');
        }
    };
</script>

<div class="max-w-[1024px] mx-auto">
    <h3 class="text-3xl">User Details</h3>
    {selectedRole?.value}
    <div class="grid grid-cols-2 gap-2 p-4 pb-2">
        <NewSensorInput
            name="username"
            label="Username"
            bind:value={user.username}
            disabled
            type="text"
            errors={fieldErrors}
        />
        <NewSensorInput
            name="name"
            label="Name"
            bind:value={user.name}
            disabled={!editing}
            type="text"
            errors={fieldErrors}
        />
        <Label
            for="type"
            class="flex items-center justify-between text-base font-semibold"
        >
            Role
            {#if fieldErrors['role']}
                <span class="text-sm font-normal italic text-red-400"
                    >{fieldErrors['type']}</span
                >
            {/if}
        </Label>
        <Select.Root
            bind:selected={selectedRole}
            onSelectedChange={(s) => {
                if (s) {
                    user.role = s.value;
                }
            }}
            disabled={props.currentUser.role != 'admin' || !editing}
            name="type"
        >
            <Select.Trigger
                class={fieldErrors['role'] ? 'border-2 border-red-600' : ''}
            >
                <Select.Value />
            </Select.Trigger>
            <Select.Content>
                {#each ROLES as { label, value }}
                    <Select.Item {value}>{label}</Select.Item>
                {/each}
            </Select.Content>
        </Select.Root>
    </div>

    {#if globalError}
        <p class="mb-1 text-center text-sm text-red-500">{globalError}</p>
    {/if}
    <div class="flex justify-end p-2 gap-4">
        <!-- TODO: DISABLE IF THIS USER IS CURRENT USER -->
        <Button
            size="bold"
            on:click={handleDelete}
            disabled={props.currentUser.id === user.id || !editing}
            variant="destructive">Delete</Button
        >
        {#if editing}
            <Button size="bold" on:click={cancelEditing} variant="outline"
                >Cancel</Button
            >
        {:else}
            <Button
                variant="outline"
                size="bold"
                on:click={() => {
                    editing = true;
                }}>Edit</Button
            >
        {/if}
        <Button size="bold" on:click={handleEdit} disabled={!editing}
            >Submit</Button
        >
    </div>
</div>
