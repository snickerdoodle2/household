<script lang="ts">
    import { ModalType } from '@/types/Modal.types';
    import { openModal } from '@/utils/Modal.utils';
    import { DotsHorizontal } from 'radix-icons-svelte';
    import { categoryStore, sensorStore, sensorValueMap } from '@/stores/Stores';
    import MyDropdown from './generic/MyDropdown.svelte';
    import MovableDropdown from './generic/MovableDropdown.svelte';
    import { submitCategoryDeletion } from '@/utils/requests/Category.requests';

    const optionsDropDown = {
        triggerRef: null as HTMLButtonElement | null,
        isOpen: false,
        menuConfig: [
            { text: 'Add Sensor', callback: () => openModal(ModalType.ADD_SENSOR, undefined) },
            { text: 'Add Category', callback: () => openModal(ModalType.ADD_CATEGORY, undefined) },
        ],
    };

    const categoryDropdown = {
        x: 0,
        y: 0,
        category: '',
        isOpen: false,
        menuConfig: [
            { text: 'Error actions not set', callback: () => console.error('actions not set!') },
        ],
    };

    function openDropdown(event: MouseEvent, category: string) {
        event.preventDefault(); // Prevent default context menu

        categoryDropdown.x = event.clientX;
        categoryDropdown.y = event.clientY;
        categoryDropdown.category = category;
        categoryDropdown.isOpen = true;
        categoryDropdown.menuConfig = [
            {
                text: 'Modify Category',
                callback: () => openModal(ModalType.MODIFY_CATEGORY, category),
            },
            {
                text: 'Delete Category',
                callback: () =>
                    openModal(ModalType.CONFIRMATION_MODAL, {
                        message: 'Are you sure you want to proceed?',
                        acceptText: 'Accept',
                        declineText: 'Decline',
                        onAccept: () => {
                            submitCategoryDeletion(category);
                        },
                        onDecline: () => {},
                    }),
            },
        ];
    }
</script>

<main>
    <div class="card flex items-center justify-between h-24 p-5 w-full rounded-lg p-2">
        <div class="flex flex-nowrap flex-row space-x-4 overflow-x-auto">
            {#each $categoryStore as category}
                <button
                    class="btn-primary text-xl"
                    on:contextmenu={(event) => {
                        openDropdown(event, category);
                        // event.preventDefault(); // Prevent default context menu
                        // openModal(ModalType.MODIFY_CATEGORY, category);
                    }}
                >
                    {category}
                </button>
            {/each}
        </div>

        <div>
            <button
                on:click={() => (optionsDropDown.isOpen = !optionsDropDown.isOpen)}
                bind:this={optionsDropDown.triggerRef}
                class="btn-primary"
            >
                <DotsHorizontal class="w-6 h-6 mr-2 text-white dark:text-white" />
            </button>

            <div class="absolute right-0 z-50">
                <MyDropdown
                    bind:isOpen={optionsDropDown.isOpen}
                    optionsWithCallbacks={optionsDropDown.menuConfig}
                    triggerButtonRef={optionsDropDown.triggerRef}
                />
            </div>
        </div>
    </div>

    <div class="flex flex-wrap gap-8 overflow-y-auto p-8">
        {#each $sensorStore as sensor}
            <button
                class="btn-secondary flex px-8 py-6 rounded-full"
                on:click={() => {
                    openModal(ModalType.SENSOR_DETAILS_MODAL, { sensor });
                }}
            >
                <p class="text-5xl">🌐</p>
                <div class="flex flex-col px-2">
                    <p class="text-2xl">{sensor.name}</p>
                    <p class="text-md">
                        {$sensorValueMap.find((sensorValue) => sensorValue.id === sensor.id)?.val ??
                            '---'}
                    </p>
                </div>
            </button>
        {/each}
    </div>

    <MovableDropdown
        x={categoryDropdown.x}
        y={categoryDropdown.y}
        items={categoryDropdown.menuConfig}
        bind:isOpen={categoryDropdown.isOpen}
    />
</main>
