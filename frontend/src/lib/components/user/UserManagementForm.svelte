<script lang="ts">
    import { newUserSchema } from "@/types/user";
    import NewSensorInput from '$lib/components/FormInput.svelte';
    import Label from '../ui/label/label.svelte';
    import * as Select from '$lib/components/ui/select';
    import Button from "../ui/button/button.svelte";
    import { authFetch } from "@/helpers/fetch";
    
    export let user: {
        id: string,
        username: string,
        name: string,
    } = {
        id:"new",
        username: "new_username",
    name: "new_name"
    }
    ;
    export let action: "add" | "edit";
    export let open: boolean;
    
    let selectedRole = {value: "user", label:"User"};
    let password: string;
    let globalError: string | undefined = undefined;

    let fieldErrors: Record<string, string> = {}

    async function handleSubmit() {
        const { data, success, error } = newUserSchema.safeParse({
            ...user,
            password
        });

        if (!success) {
            console.log(error.issues);
            fieldErrors = {};
            for(const issue of error.issues){
                fieldErrors[issue.path[0]] = issue.message
            }
            return;
        }

        const res = await authFetch(`/api/v1/user/${data.username}`, {
            method: action === "add" ? 'POST' : 'PUT',
            body: JSON.stringify(data),
        });

        const resJson = await res.json();
        console.log(resJson);

        if (res.ok) {
            open = false;
        } else {
            if (typeof resJson.error === 'string') {
                globalError = resJson.error;
            } else {
                fieldErrors = resJson.error;
            }
        }
    }
</script>


{#if user}
    <div>
        <h3 class="text-3xl"> {action === "add" ? "Add User" : "Edit User"} </h3>
        <div class="grid grid-cols-2 gap-2 p-4 pb-2">
            <NewSensorInput
                    name="name"
                    label="Name"
                    bind:value={user.name}
                    type="text"
                    errors={fieldErrors}
                />
            <NewSensorInput
                    name="username"
                    label="Username"
                    bind:value={user.username}
                    type="text"
                    errors={fieldErrors}
                />
            <Label
                for="type"
                class="flex items-center justify-between text-base font-semibold"
            >
                Role
                {#if fieldErrors['type']}
                    <span class="text-sm font-normal italic text-red-400"
                        >{fieldErrors['type']}</span
                    >
                {/if}
            </Label>
            <Select.Root bind:selected={selectedRole} required name="type">
                <Select.Trigger
                    class={fieldErrors['type'] ? 'border-2 border-red-600' : ''}
                >
                    <Select.Value />
                </Select.Trigger>
                <Select.Content>
                    <Select.Item value={"user"}>User</Select.Item>
                    <Select.Item value={"admin"}>Admin</Select.Item>
                </Select.Content>
            </Select.Root>
            <NewSensorInput
                    name="password"
                    label="Password"
                    bind:value={password}
                    type="text"
                    errors={fieldErrors}
                />
        </div>

        {#if globalError}
            <p class="mb-1 text-sm text-red-500 text-center">{globalError}</p>
        {/if}
        <div class ="p-2 flex justify-end w-full">
            <Button size="bold" on:click={handleSubmit}>Submit</Button>
        </div>
    </div>    
{/if}


