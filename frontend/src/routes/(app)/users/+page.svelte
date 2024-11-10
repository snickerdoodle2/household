<script lang="ts">
    import type { User } from "@/types/user";
    import type { PageData } from "./$types";
    import { onMount } from "svelte";
    import Button from "@/components/ui/button/button.svelte";
    import { Trash, Pencil1 } from 'svelte-radix';
    import { authFetch } from "@/helpers/fetch";
    export let data: PageData;
    import * as Dialog from '$lib/components/ui/dialog';
    import UserManagementForm from "@/components/user/UserManagementForm.svelte";
    
    let users: User[] = [];

    let modalOpen = false;
    let modalData: User | null = null 

    onMount(async () => {
        users = await data.users
    })

    const handleUserAdd = () => {
        modalData = null;
        modalOpen = true;
    }

    const handleUserEdit = (user: User) => {
        modalData = user;
        modalOpen = true;
    }

    const handleDelete = async (username: string) => {
        // TODO: ask for confirmation!!!
        const res = await authFetch(`/api/v1/user/${username}`, {
            method: 'DELETE',
        });

        if (res.ok) {
            console.log(await res.json());
        } else {
            const resJson = await res.json();
            console.log(resJson);
        }
    };


</script>


<div class="flex flex-col w-full justify-start p-8">
    <div class="w-full flex justify-between">
        <h2 class="text-3xl font-bold pb-5">All Users ({users.length})</h2>
        <Button class="mt-auto" variant="outline" size="sm" on:click={handleUserAdd}> Add user </Button>    
    </div>
    <div class="grid grid-cols-[3fr_3fr_3fr_1fr] w-full gap-2 items-center">
        <p class="pb-1 text-xl font-bold">Username</p>
        <p class="pb-1 text-xl font-bold">Role</p>
        <p class="pb-1 text-xl font-bold">Created at</p>
        <div></div>

        {#each users as user}
            <div class="border-t-1 border-gray-300">
                <h3 class="text-md font-bold">{user.username} <span class="text-sm font-thin">({user.name})</span></h3>
                <p class="text-xs font-thin">{user.id}</p>
            </div>

            <div class="items-start border-gray-300">
                <div class="inline-block border-2 rounded-xl px-2 font-bold text-sm {Math.random() > 0.5 ? "bg-orange-300 border-orange-500 text-orange-500" : "bg-blue-300 border-blue-500 text-blue-500"} ">
                    user/admin
                </div>
            </div>

            <p class="font-thin">{user.created_at.toLocaleString('en-US', { 
                    hour: 'numeric', minute: 'numeric', second: 'numeric', 
                    hour12: false 
                })} - {user.created_at.toLocaleString('en-GB', { 
                    day: 'numeric', month: 'long', year: 'numeric' 
                })}
            </p>

            <div class="flex justify-end items-center">
                <Button variant="outline" size="icon" on:click={() => handleUserEdit(user)}> 
                    <Pencil1 class="w-5 h-5" /> 
                </Button>    
                <Button variant="outline" size="icon" on:click={() => handleDelete(user.username)}> 
                    <Trash class="w-5 h-5"/>
                </Button>    
            </div>
        {/each}
    </div>
</div>

<Dialog.Root
    bind:open={modalOpen}
    onOpenChange={(opened) => {
        if (!opened) modalData = null;
        console.log("openedChange")
    }}
>
    <Dialog.Portal>
        <Dialog.Overlay />
        <Dialog.Content
            class="flex max-w-none items-center justify-center px-8 py-4 md:w-fit"
        >
            {#if modalData}
                <UserManagementForm user={modalData} action ="edit" bind:open={modalOpen}/>
            {:else}
                <UserManagementForm action ="add" bind:open={modalOpen}/>
            {/if}
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>
