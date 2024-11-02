<script lang="ts">
    import type { SequenceAction } from "@/types/sequence";
    import Action from "./Action.svelte";
    import Button from "../ui/button/button.svelte";
    import { Trash } from "radix-icons-svelte";
    import { Plus } from "svelte-radix";

    export let actions: SequenceAction[];
    export let sensors: { label: string; value: string }[]
    export let editing: boolean = true;

    function addAction(){
        actions.push({
            target: sensors[0].value,
            value: 0, 
            msDelay: 0
        })
        actions = actions // to trigger reactive logic
    }

    function removeAction(idx: number){
        actions.splice(idx, 1)
        actions = actions // to trigger reactive logic
    }

    $: console.log(actions, "from builder")
</script>

<div class="w-full min-w-[40rem]">
    <ul>
        {#each actions as action, idx}
            <li class="flex p-1">
                <Action bind:sensors bind:action bind:editing />
                {#if editing}
                    <Button
                        on:click={() => removeAction(idx)}
                        variant="outline"
                        size="sm"
                        class="ml-2"
                    >
                        <Trash class="w-4" />
                    </Button>
                {/if}
            </li>
        {/each}
        {#if editing}
            <li class="flex items-center justify-center p-2">
                <Button on:click={addAction} variant="outline" size="sm">
                    <Plus class="w-4" />
                </Button>
            </li>
        {/if}
    </ul>
</div>
