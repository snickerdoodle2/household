<script>
    import { ModalType } from '@/types/Modal.types';
    import { openModal } from '@/utils/Modal.utils';
    import { DotsHorizontal } from 'radix-icons-svelte';
    import { Button } from 'flowbite-svelte';
    import { sensorStore, sensorValueMap } from '@/stores/Stores';

    const categories = [
        'Kitchen',
        'Living Room',
        'Bedroom',
        'Bathroom',
        'Garage',
        'Garden',
    ];
</script>

<main>
    <div
        class="flex items-center justify-between bg-black w-full h-0.3 rounded-lg p-2"
    >
        <div class="flex flex-nowrap flex-row space-x-4 overflow-x-auto">
            {#each categories as category}
                <Button class="bg-secondary rounded-lg">
                    {category}
                </Button>
            {/each}
        </div>

        <Button
            on:click={() => openModal(ModalType.ADD_SENSOR, undefined)}
            class="bg-secondary rounded-lg"
        >
            <DotsHorizontal class="w-6 h-6 mr-2 text-white dark:text-white" />
        </Button>
    </div>

    {#each $sensorStore as sensor}
        <div class="flex flex-wrap overflow-y-auto p-2">
            <button class="flex px-3 py-4 bg-black rounded-full">
                <p class="text-4xl">🌐</p>
                <div class="flex flex-col px-2">
                    <p class="text-md">{sensor.name}</p>
                    <p class="text-sm">
                        {$sensorValueMap.find(
                            (sensorValue) => sensorValue.id === sensor.id
                        )?.val ?? '---'}
                    </p>
                </div>
            </button>
            <!-- <h1 class="text-4xl font-bold text-white">Sensor Page</h1> -->
        </div>
    {/each}
</main>
