<script lang="ts">
    import { currentPageStore } from '@/stores/Stores';
    import { PageType } from '@/types/Page.types';
    import { Dropdown, DropdownItem, Navbar, NavBrand } from 'flowbite-svelte';
    import { ChevronDownOutline, ListOutline } from 'flowbite-svelte-icons';
    import SensorDisplay from './SensorDisplay.svelte';
    import { onDestroy, onMount } from 'svelte';
    import { initializeSensorData, syncSensorValues } from '@/utils/Sync.utils';
    import { SENSOR_VALUE_INTERVAL } from '@/config/const';
    import Button from './ui/button/button.svelte';

    let syncInterval: number;

    onMount(() => {
        // Initialization
        initializeSensorData();

        // Sync interval
        syncInterval = setInterval(() => {
            syncSensorValues();
        }, SENSOR_VALUE_INTERVAL.toMilliseconds());
    });

    onDestroy(() => {
        clearInterval(syncInterval);
    });
</script>

<main class="bg-background text-foreground px-[2.5vw] py-[2.5vh]">
    <div class="card-muted flex rounded-lg h-[10vh] md-[2.5vh] p-10">
        <div class="flex items-center">
            <button id="list-outline">
                <ListOutline
                    class="w-12 h-12 mr-2 text-foreground"
                />
            </button>
            <span class="font-semibold text-3xl">Nazwa Naszego Systemu</span>
        </div>

        <Dropdown
            triggeredBy="#list-outline"
            placement="bottom"
            class="bg-popover rounded-lg"
        >
            {#each ['See Rules', 'Set Server URL', 'Do ---- nothing', 'just --------- placeholder'] as option}
                <DropdownItem>
                    <p class="text-xl">{option}</p>
                </DropdownItem>
            {/each}
        </Dropdown>
    </div>

    <div class="card-muted p-0 rounded-lg h-[80vh] my-[2.5vw]">
        {#if $currentPageStore === PageType.SENSOR}
            <SensorDisplay />
        {/if}
    </div>
</main>

<style>
</style>
