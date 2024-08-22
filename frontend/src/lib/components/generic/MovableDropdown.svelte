<!-- src/components/DropdownMenu.svelte -->
<script lang="ts">
    export let x: number;
    export let y: number;
    export let items: { text: string; callback: () => void }[] = [];
    export let isOpen: boolean;

    let isMouseOverDropdown = false;
    let destroyTimeout = 0;

    function handleMouseEnter() {
        isMouseOverDropdown = true;
        clearTimeout(destroyTimeout);
    }

    function handleMouseLeave() {
        isOpen = false;
        isMouseOverDropdown = false;
    }

    function handleDropdownClick() {
        if (isMouseOverDropdown) {
            console.log('Mouse is over the dropdown.');
        }
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === 'Enter' || event.key === ' ') {
            handleDropdownClick();
        }
    }

    // Reactively watch isOpen and set a timeout to close the dropdown if not hovered
    $: if (isOpen) {
        destroyTimeout = setTimeout(() => {
            if (!isMouseOverDropdown) {
                isOpen = false;
            }
        }, 1000);

        // Clear timeout if dropdown is closed
        if (!isOpen) {
            clearTimeout(destroyTimeout);
        }
    }
</script>

{#if isOpen}
    <div
        class="dropdown-menu"
        role="button"
        tabindex="0"
        style="top: {y}px; left: {x}px;"
        on:mouseenter={handleMouseEnter}
        on:mouseleave={handleMouseLeave}
        on:click={handleDropdownClick}
        on:keydown={handleKeydown}
    >
        {#each items as item}
            <button
                on:click={() => {
                    item.callback();
                    isOpen = false;
                }}
                class="dropdown-item"
            >
                {item.text}
            </button>
        {/each}
    </div>
{/if}

<style>
    .dropdown-menu {
        position: absolute;
        background-color: hsl(var(--background));
        border: 1px solid hsl(var(--border));
        border-radius: var(--radius);
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
        z-index: 1000;
        padding: 0.5rem;
    }

    .dropdown-item {
        display: block;
        padding: 0.5rem 1rem;
        border-radius: var(--radius);
        color: hsl(var(--foreground));
        background: transparent;
        border: none;
        text-align: left;
        cursor: pointer;
        width: 100%;
        transition: background-color 0.2s;
    }

    .dropdown-item:hover {
        background-color: hsl(var(--primary-hover));
    }
</style>
