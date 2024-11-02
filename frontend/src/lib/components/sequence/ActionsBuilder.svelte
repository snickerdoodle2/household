<script lang="ts">
    import type { SequenceAction } from "@/types/sequence";
    import Action from "./Action.svelte";
    import Button from "../ui/button/button.svelte";
    import { Plus, ArrowDown, Trash} from "svelte-radix";

    export let actions: SequenceAction[];
    export let sensors: { label: string; value: string }[]
    export let editing: boolean = true;
    export let fieldErrors: Record<number, string[]> = []

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

    function moveDown(idx: number){
        if (idx >= actions.length - 1) return;
        const tmp = actions[idx]  
        actions[idx] = actions[idx + 1]
        actions[idx + 1] = tmp
    }
</script>

<div class="w-full min-w-[40rem]">
    <ul>
        {#each actions as action, idx}
            <li class="flex p-1 items-center justify-center">
                <Action
                    bind:sensors
                    bind:action
                    bind:editing
                    bind:errorFields={fieldErrors[idx]}
                />
                {#if editing}
                    <Button
                        on:click={() => moveDown(idx)}
                        variant="outline"
                        class="ml-5"
                        size="sm"
                    >
                        <ArrowDown class="w-4" />
                    </Button>
                    <Button
                        on:click={() => removeAction(idx)}
                        variant="outline"
                        class="ml-1"
                        size="sm"
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
