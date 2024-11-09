<script lang="ts">
    import type { User } from "@/types/user";
    import type { PageData } from "./$types";
    import { onMount } from "svelte";
    import Button from "@/components/ui/button/button.svelte";
    import { Trash, Pencil1 } from 'svelte-radix';
export let data: PageData;
let users: User[] = [];

onMount( async () => {
    users = await data.users
})
</script>


<div class="flex flex-col w-full justify-start p-8">
    <div class="w-full flex justify-between">
        <h2 class="text-3xl font-bold pb-5">All Users ({users.length})</h2>
        <Button class="mt-auto" variant="outline" size="sm"> Add user </Button>    
    </div>
    <div class="grid grid-cols-[3fr_3fr_3fr_1fr] w-full gap-2 items-center">
        <p class="pb-1 text-xl font-bold">Username</p>
        <p class="pb-1 text-xl font-bold">Role</p>
        <p class="pb-1 text-xl font-bold">Created at</p>
        <div></div>

        {#each users as user}
            <div class="border-t-1 border-gray-300">
                <h3 class="text-md font-bold">{user.username}</h3>
                <p class="text-sm font-thin">{user.name}</p>
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
                <Button variant="outline" size="icon"> 
                    <Pencil1 class="w-5 h-5"/> 
                </Button>    
                <Button variant="outline" size="icon"> 
                    <Trash class="w-5 h-5"/>
                </Button>    
            </div>
        

        {/each}
    
        
    </div>

</div>