<script lang="ts">
    import { currentPageStore } from '@/stores/Stores';
    import { PageType } from '@/types/Page.types';
    import { ListOutline } from 'flowbite-svelte-icons';
    import SensorDisplay from './SensorDisplay.svelte';
    import { onDestroy, onMount } from 'svelte';
    import { initializeStores, syncSensorValues } from '@/utils/Sync.utils';
    import { SENSOR_VALUE_INTERVAL } from '@/config/const';
    import MyDropdown from './generic/MyDropdown.svelte';
    import RulesDisplay from './RulesDisplay.svelte';
    import { openPage } from '@/utils/Page.utils';

    let syncInterval: number;
    let dropDown = {
        open: false,
        trigger: null as HTMLButtonElement | null,
    };

    const staticDropDownMenuActions = [
        {
            text: 'Set Server URL',
            callback: () => console.log('Set Server URL clicked'),
        },
        {
            text: 'Do ---- nothing',
            callback: () => console.log('Do ---- nothing clicked'),
        },
        {
            text: 'just --------- placeholder',
            callback: () => console.log('Placeholder clicked'),
        },
    ];

    $: dropDownMenuConfig =
        $currentPageStore === PageType.SENSOR
            ? [
                  ...staticDropDownMenuActions,
                  { text: 'See Rules', callback: () => openPage(PageType.RULE) },
              ]
            : [
                  ...staticDropDownMenuActions,
                  { text: 'See Sensors', callback: () => openPage(PageType.SENSOR) },
              ];

    onMount(() => {
        // Initialization
        initializeStores();

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
            <button
                id="list-outline"
                bind:this={dropDown.trigger}
                on:click={() => (dropDown.open = !dropDown.open)}
            >
                <ListOutline class="w-12 h-12 mr-2 text-foreground" />
            </button>
            <span class="font-semibold text-3xl">Nazwa Naszego Systemu</span>
        </div>

        <div class="absolute mt-8">
            <MyDropdown
                optionsWithCallbacks={dropDownMenuConfig}
                triggerButtonRef={dropDown.trigger}
                bind:isOpen={dropDown.open}
            ></MyDropdown>
        </div>
    </div>

    <div class="card-muted p-0 rounded-lg h-[80vh] my-[2.5vw]">
        {#if $currentPageStore === PageType.SENSOR}
            <SensorDisplay />
        {:else if $currentPageStore === PageType.RULE}
            <RulesDisplay />
        {/if}
    </div>
</main>

<style>
</style>
