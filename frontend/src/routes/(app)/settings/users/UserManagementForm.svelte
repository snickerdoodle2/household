<script lang="ts">
    import { userSchema, type User } from '@/types/user';
    import NewSensorInput from '$lib/components/FormInput.svelte';
    import Button from '$lib/components/ui/button/button.svelte';
    import { authFetch } from '@/helpers/fetch';

    type Props =
        | {
              user: User;
              action: 'edit';
          }
        | {
              user: undefined;
              action: 'add';
          };

    const props: Props = $props();

    let user = $state(
        props.user ?? {
            id: '',
            username: '',
            name: '',
        }
    );

    const cancelEditing = () => {
        user = props.user ?? {
            id: '',
            username: '',
            name: '',
        };
        editing = false;
    };

    let editing = $state(false);

    let globalError: string | undefined = $state();
    let fieldErrors: Record<string, string> = $state({});

    async function handleEdit() {
        const { data, success, error } = userSchema.safeParse({
            id: user.id,
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
        console.log(resJson);

        if (typeof resJson.error === 'string') {
            globalError = resJson.error;
        } else {
            fieldErrors = resJson.error;
        }
    }
</script>

<div>
    <h3 class="text-3xl">
        {props.action === 'add' ? 'Add User' : 'User Details'}
    </h3>
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
    {#if props.action === 'edit'}
        <div class="flex justify-end p-2 gap-4">
            <!-- TODO: DISABLE IF THIS USER IS CURRENT USER -->
            <Button
                size="bold"
                on:click={() => {}}
                disabled={!editing}
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
            <Button
                size="bold"
                on:click={props.action === 'edit' ? handleEdit : () => {}}
                >Submit</Button
            >
        </div>
    {:else}{/if}
</div>
