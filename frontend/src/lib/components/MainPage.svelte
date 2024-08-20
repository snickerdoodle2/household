<script lang="ts">
    import { currentPageStore } from '@/stores/Stores';
    import { PageType } from '@/types/Page.types';
    import { Dropdown, DropdownItem, Navbar, NavBrand } from 'flowbite-svelte';
    import { ListOutline } from 'flowbite-svelte-icons';
    import SensorDisplay from './SensorDisplay.svelte';
    import { onDestroy, onMount } from 'svelte';
    import { initializeSensorData, syncSensorValues } from '@/utils/Sync.utils';
    import { SENSOR_VALUE_INTERVAL } from '@/config/const';

    let syncInterval: number;

    onMount(() => {
        // Initialization
        initializeSensorData()

        // Sync interval
        syncInterval = setInterval(() => {
            syncSensorValues();
        }, SENSOR_VALUE_INTERVAL.toMilliseconds());
    });

    onDestroy(() => {
        clearInterval(syncInterval);
    });
</script>

<main class="bg-primary px-[2.5vw] py-[2.5vh]">
    <Navbar class="bg-card flex rounded-lg h-[10vh] md-[2.5vh]">
        <NavBrand href="/">
            <div class="flex items-center">
                <button
                    class="flex items-center bg-transparent border-none cursor-pointer p-0"
                >
                    <ListOutline
                        class="w-6 h-6 mr-2 text-white dark:text-white"
                    />
                    <span class="font-semibold text-xl"
                        >Nazwa Naszego Systemu</span
                    >
                </button>
                <Dropdown
                    class="bg-popover absolute mt-2 w-[5hv] right-0 top-1/2 rounded-lg"
                >
                    <DropdownItem>Dashboard</DropdownItem>
                    <DropdownItem>Settings</DropdownItem>
                    <DropdownItem>Earnings</DropdownItem>
                    <DropdownItem>Sign out</DropdownItem>
                </Dropdown>
            </div>
        </NavBrand>
    </Navbar>

    <div class="bg-card rounded-lg h-[80vh] my-[2.5vw]">
        {#if $currentPageStore === PageType.SENSOR}
            <SensorDisplay />
        {/if}
    </div>
</main>
