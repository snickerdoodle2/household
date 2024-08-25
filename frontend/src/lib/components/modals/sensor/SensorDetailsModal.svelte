<!-- src/routes/SensorDetailsModal.svelte -->
<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { SensorType, type Sensor, type SensorData } from '@/types/Sensor.types';
    import { validateName } from '@/utils/Misc.utils';
    import { closeModal, isModalData } from '@/utils/Modal.utils';
    import { ModalType } from '@/types/Modal.types';
    import { openedModalStore } from '@/stores/Stores';
    import { get } from 'svelte/store';
    import SensorInputModal from './SensorInputModal.svelte';
    let storeData = get(openedModalStore);
    let sensor = getSensorData();

    function getSensorData() {
        if (storeData && isModalData(ModalType.SENSOR_DETAILS_MODAL, storeData)) {
            return storeData.data.sensor;
        } else {
            closeModal();
            const sensor: Sensor = {
                id: '',
                name: '',
                type: SensorType.BINARY_SENSOR,
                uri: '',
                created_at: '',
                version: 0,
                refresh_rate: 0,
            };
            return sensor;
        }
    }

    let isEditing = false;

    const handleEditSubmit = async (data: SensorData) => {
        sensor.name = data.name;
        sensor.refresh_rate = data.refresh_rate;
        sensor.uri = data.uri;
        sensor.type = data.type as SensorType;
    };

    // Handlers for other buttons
    const handleRemove = () => {
        console.log(`Remove sensor with ID: ${sensor.id}`);
    };

    const handleMonitor = () => {
        console.log(`Monitor sensor with ID: ${sensor.id}`);
    };

    const handleSeeRules = () => {
        console.log(`See rules for sensor with ID: ${sensor.id}`);
    };
</script>

<main>
    <div class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
        <div class="bg-background rounded-lg shadow-lg w-full max-w-lg p-6 relative">
            <button
                type="button"
                class="absolute top-2 right-2 btn-exit"
                on:click={closeModal}
                on:close={() => (isEditing = false)}
            >
                &times;
            </button>
            <h2 class="text-2xl font-bold mb-4">{isEditing ? 'Edit Sensor' : 'Sensor Details'}</h2>

            {#if isEditing}
                <SensorInputModal
                    title={'Edit Sensor'}
                    sensorData={{ ...sensor }}
                    onSubmit={handleEditSubmit}
                />
            {:else}
                <ul class="space-y-2">
                    <li><strong>ID:</strong> {sensor.id}</li>
                    <li><strong>Name:</strong> {sensor.name}</li>
                    <li><strong>Type:</strong> {sensor.type}</li>
                    <li><strong>URI:</strong> {sensor.uri}</li>
                    <li><strong>Created At:</strong> {sensor.created_at}</li>
                    <li><strong>Version:</strong> {sensor.version}</li>
                    <li><strong>Refresh Rate:</strong> {sensor.refresh_rate}</li>
                </ul>

                <div class="mt-6 flex justify-end gap-4">
                    <button
                        class="btn-primary"
                        on:click={() => {
                            isEditing = true;
                        }}
                    >
                        Edit
                    </button>
                    <button class="btn-danger" on:click={handleRemove}> Remove </button>
                    <button class="btn-secondary" on:click={handleMonitor}> Monitor </button>
                    <button class="btn-info" on:click={handleSeeRules}> See Rules </button>
                </div>
            {/if}
        </div>
    </div>
</main>

<style>
    .invalid {
        border-color: red;
        border-style: solid;
        border-width: 2px;
        border-radius: var(--radius);
    }
</style>
