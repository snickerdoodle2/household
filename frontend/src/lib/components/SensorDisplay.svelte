<script lang="ts">
    import { ModalType } from '@/types/Modal.types';
    import { openModal } from '@/utils/Modal.utils';
    import { DotsHorizontal } from 'radix-icons-svelte';
    import { categoryStore, sensorStore, sensorValueMap } from '@/stores/Stores';
    import MyDropdown from './generic/MyDropdown.svelte';

    const dropDown = {
        triggerRef: null as HTMLButtonElement | null,
        isOpen: false,
        menuConfig: [
            { text: 'Add Sensor', callback: () => openModal(ModalType.ADD_SENSOR, undefined) },
            {
                text: 'Add Category',
                callback: () => console.log('open add category modal'),
            },
        ],
    };
</script>

<main>
    <div class="card flex items-center justify-between h-24 p-5 w-full rounded-lg p-2">
        <div class="flex flex-nowrap flex-row space-x-4 overflow-x-auto">
            {#each $categoryStore as category}
                <button class="btn-primary text-2xl">
                    {category}
                </button>
            {/each}
        </div>

        <div>
            <button
                on:click={() => (dropDown.isOpen = !dropDown.isOpen)}
                bind:this={dropDown.triggerRef}
                class="btn-primary"
            >
                <DotsHorizontal class="w-6 h-6 mr-2 text-white dark:text-white" />
            </button>

            <div class="absolute right-0">
                <MyDropdown
                    bind:isOpen={dropDown.isOpen}
                    optionsWithCallbacks={dropDown.menuConfig}
                    triggerButtonRef={dropDown.triggerRef}
                />
            </div>
        </div>
    </div>

    {#each $sensorStore as sensor}
        <div class="flex flex-wrap overflow-y-auto p-8">
            <button class="btn-secondary flex px-8 py-6 rounded-full">
                <p class="text-5xl">🌐</p>
                <div class="flex flex-col px-2">
                    <p class="text-2xl">{sensor.name}</p>
                    <p class="text-md">
                        {$sensorValueMap.find((sensorValue) => sensorValue.id === sensor.id)?.val ??
                            '---'}
                    </p>
                </div>
            </button>
            <!-- <h1 class="text-4xl font-bold text-white">Sensor Page</h1> -->
        </div>
    {/each}
</main>
