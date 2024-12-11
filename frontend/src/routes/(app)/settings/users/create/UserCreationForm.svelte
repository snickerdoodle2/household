<script lang="ts">
    import { newUserSchema, type NewUser } from '@/types/user';
    import NewSensorInput from '$lib/components/FormInput.svelte';
    import Button from '$lib/components/ui/button/button.svelte';
    import { authFetch } from '@/helpers/fetch';
    import { goto } from '$app/navigation';

    let user: NewUser = $state({
        name: '',
        username: '',
        password: '',
        confirmPassword: '',
    });

    let globalError: string | undefined = $state();
    let fieldErrors: Record<string, string> = $state({});

    async function handleCreate() {
        const { success, error, data } = newUserSchema.safeParse(user);
        if (!success) {
            console.log(error.issues);
            fieldErrors = {};
            for (const issue of error.issues) {
                fieldErrors[issue.path[0]] = issue.message;
            }
            return;
        }

        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        const { confirmPassword, ...userData } = data;

        const res = await authFetch('/api/v1/user', {
            method: 'POST',
            body: JSON.stringify(userData),
        });

        const body = await res.json();
        console.log(body);

        if (res.ok) {
            goto('/settings/users');
        }
    }
</script>

<div class="max-w-[1024px] mx-auto">
    <h3 class="text-3xl">User Details</h3>
    <div class="grid grid-cols-2 gap-2 p-4 pb-2">
        <NewSensorInput
            name="username"
            label="Username"
            bind:value={user.username}
            type="text"
            errors={fieldErrors}
        />
        <NewSensorInput
            name="name"
            label="Name"
            bind:value={user.name}
            type="text"
            errors={fieldErrors}
        />
        <NewSensorInput
            name="password"
            label="Password"
            bind:value={user.password}
            type="password"
            errors={fieldErrors}
        />
        <NewSensorInput
            name="confirmPassword"
            label="Confirm Password"
            bind:value={user.confirmPassword}
            type="password"
            errors={fieldErrors}
        />
        <!-- <Label -->
        <!--     for="type" -->
        <!--     class="flex items-center justify-between text-base font-semibold" -->
        <!-- > -->
        <!--     Role -->
        <!--     {#if fieldErrors['type']} -->
        <!--         <span class="text-sm font-normal italic text-red-400" -->
        <!--             >{fieldErrors['type']}</span -->
        <!--         > -->
        <!--     {/if} -->
        <!-- </Label> -->
        <!-- <Select.Root bind:selected={selectedRole} required name="type"> -->
        <!--     <Select.Trigger -->
        <!--         class={fieldErrors['type'] ? 'border-2 border-red-600' : ''} -->
        <!--     > -->
        <!--         <Select.Value /> -->
        <!--     </Select.Trigger> -->
        <!--     <Select.Content> -->
        <!--         <Select.Item value={'user'}>User</Select.Item> -->
        <!--         <Select.Item value={'admin'}>Admin</Select.Item> -->
        <!--     </Select.Content> -->
        <!-- </Select.Root> -->
    </div>

    {#if globalError}
        <p class="mb-1 text-center text-sm text-red-500">{globalError}</p>
    {/if}
    <div class="flex justify-end p-2 gap-4">
        <Button size="bold" on:click={handleCreate}>Submit</Button>
    </div>
</div>
