<script lang="ts">
    import { onMount, onDestroy } from 'svelte';

    export let isOpen = false;
    export let optionsWithCallbacks: { text: string; callback: () => void }[] = [];
    export let triggerButtonRef: HTMLButtonElement | null = null; // Define the type for dropdownRef

    let dropdownRef: HTMLDivElement | null = null; // Define the type for dropdownRef

    function handleClickOutside(event: MouseEvent) {
        if (
            isOpen &&
            dropdownRef &&
            !dropdownRef.contains(event.target as Node) &&
            triggerButtonRef &&
            !triggerButtonRef.contains(event.target as Node)
        ) {
            isOpen = false;
        }
    }

    onMount(() => {
        document.addEventListener('click', handleClickOutside);
    });

    onDestroy(() => {
        document.removeEventListener('click', handleClickOutside);
    });
</script>

{#if isOpen}
    <div bind:this={dropdownRef} class="flex flex-col bg-popover rounded-lg mt-2 shadow-lg w-64">
        {#each optionsWithCallbacks as option}
            <button
                class="p-2 hover:bg-primary-hover rounded-lg text-xl"
                on:click={() => {
                    option.callback();
                    isOpen = false;
                }}
            >
                {option.text}
            </button>
        {/each}
        <!-- <ul class="p-2">
            <li class="p-2 hover:bg-primary-hover">Option 1</li>
            <li class="p-2 hover:bg-primary-hover">Option 2</li>
            <li class="p-2 hover:bg-primary-hover">Option 3</li>
        </ul> -->
    </div>
{/if}
