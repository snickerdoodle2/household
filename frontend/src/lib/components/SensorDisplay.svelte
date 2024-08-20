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
        class="card flex items-center justify-between h-24 p-5 w-full  rounded-lg p-2"
    >
        <div class="flex flex-nowrap flex-row space-x-4 overflow-x-auto">
            {#each categories as category}
                <Button class="btn-primary rounded-lg text-2xl ">
                    {category}
                </Button>
            {/each}
        </div>

        <Button
            on:click={() => openModal(ModalType.ADD_SENSOR, undefined)}
            class="btn-primary rounded-lg"
        >
            <DotsHorizontal class="w-6 h-6 mr-2 text-white dark:text-white" />
        </Button>
    </div>

    {#each $sensorStore as sensor}
        <div class="flex flex-wrap overflow-y-auto p-8">
            <button class="btn-secondary flex px-8 py-6 rounded-full">
                <p class="text-5xl">🌐</p>
                <div class="flex flex-col px-2">
                    <p class="text-2xl">{sensor.name}</p>
                    <p class="text-md">
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
