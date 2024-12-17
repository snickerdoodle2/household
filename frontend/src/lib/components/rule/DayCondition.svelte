<script lang="ts">
    import type { RuleDayType } from '@/types/rule';
    import { onMount } from 'svelte';
    import Button from '../ui/button/button.svelte';
    import { Disc, Pencil1 } from 'radix-icons-svelte';
    import DayConditionSelect from './DayConditionSelect.svelte';

    type Props = {
        internal: RuleDayType;
        editingDisabled?: boolean;
        children?: import('svelte').Snippet;
    };

    let {
        internal = $bindable(),
        editingDisabled = $bindable(false),
        children,
    }: Props = $props();

    let editing = $state(false);

    let dropDownsOpen = $state(false);
    let format = $state('');
    let forcedFormat = $state('');

    let wrappingDiv: HTMLDivElement | undefined = $state();

    function toggleEditing() {
        editing = !editing;
    }

    function save() {
        if (!format) return false;
        internal.format = format;
        forcedFormat = format;
        return true;
    }

    function syncInternalValues() {
        forcedFormat = internal.format;
    }

    function handleClick(event: MouseEvent) {
        if (
            wrappingDiv &&
            !wrappingDiv.contains(event.target as Node) &&
            !event.defaultPrevented &&
            !Object.values(dropDownsOpen).some((open) => open)
        ) {
            editing = false;
            syncInternalValues();
        }
    }

    onMount(() => {
        // initialize values
        syncInternalValues();

        // Capture phase to ensure clicks inside are checked before propagation finishes
        document.addEventListener('pointerdown', handleClick, true);
    });
</script>

<div
    class="flex w-full flex-row items-center gap-2 whitespace-nowrap"
    bind:this={wrappingDiv}
>
    <DayConditionSelect
        bind:format
        bind:dropDownsOpen
        disabled={!editing}
        bind:forcedFormat
    ></DayConditionSelect>

    {#if !editingDisabled}
        {#if editing}
            <Button
                variant="outline"
                size="icon"
                on:click={() => {
                    if (save()) toggleEditing();
                }}
            >
                <Disc />
            </Button>
        {:else}
            <Button variant="outline" size="icon" on:click={toggleEditing}>
                <Pencil1 />
            </Button>
        {/if}
    {/if}

    {@render children?.()}
</div>

<style>
</style>
