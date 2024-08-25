<!-- src/routes/SensorDetailsModal.svelte -->
<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { SensorType, type Sensor } from '@/types/Sensor.types';
    import { validateName } from '@/utils/Misc.utils';
    import { closeModal, isModalData } from '@/utils/Modal.utils';
    import { ModalType } from '@/types/Modal.types';
    import { openedModalStore } from '@/stores/Stores';
    import { get } from 'svelte/store';
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
    let editedSensor: Sensor = { ...sensor };

    let isInvalidName: boolean = false;
    let isInvalidType: boolean = false;
    let isInvalidURI: boolean = false;
    let isInvalidRefreshRate: boolean = false;

    const dispatcher = createEventDispatcher();

    const handleEditSubmit = async (event: Event) => {
        event.preventDefault();
        if (!validateForm()) return;
        Object.assign(sensor, editedSensor);
        isEditing = false;
        dispatcher('updated', sensor);
    };

    function validateForm(): boolean {
        let isInvalid = false;
        const uriRegex =
            /^(25[0-5]|2[0-4]\d|1\d\d|\d\d?)\.(25[0-5]|2[0-4]\d|1\d\d|\d\d?)\.(25[0-5]|2[0-4]\d|1\d\d|\d\d?)\.(25[0-5]|2[0-4]\d|1\d\d|\d\d?):([1-9]\d{0,4})$/;

        if (validateName(editedSensor.name).isError) {
            isInvalidName = true;
            isInvalid = true;
        } else {
            isInvalidName = false;
        }

        if (Object.values(SensorType).indexOf(editedSensor.type) === -1) {
            isInvalidType = true;
            isInvalid = true;
        } else {
            isInvalidType = false;
        }

        if (!uriRegex.test(editedSensor.uri)) {
            isInvalidURI = true;
            isInvalid = true;
        } else {
            isInvalidURI = false;
        }

        if (editedSensor.refresh_rate <= 0) {
            isInvalidRefreshRate = true;
            isInvalid = true;
        } else {
            isInvalidRefreshRate = false;
        }

        return !isInvalid;
    }

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
            <button type="button" class="absolute top-2 right-2 btn-exit" on:click={closeModal}>
                &times;
            </button>
            <h2 class="text-2xl font-bold mb-4">{isEditing ? 'Edit Sensor' : 'Sensor Details'}</h2>

            {#if isEditing}
                <form on:submit={handleEditSubmit} class="space-y-4">
                    <div>
                        <label for="name" class="block text-sm font-medium mb-1">Name:</label>
                        <input
                            type="text"
                            id="name"
                            bind:value={editedSensor.name}
                            class={`input-field w-full ${isInvalidName ? 'border-red-500' : ''}`}
                        />
                    </div>

                    <div>
                        <label for="uri" class="block text-sm font-medium mb-1">URI:</label>
                        <input
                            type="text"
                            id="uri"
                            bind:value={editedSensor.uri}
                            class={`input-field w-full ${isInvalidURI ? 'border-red-500' : ''}`}
                        />
                    </div>

                    <!-- Type and Refresh Rate Fields -->
                    <div class="flex gap-4">
                        <!-- Sensor Type Field -->
                        <div class={`w-3/4 ${isInvalidType ? 'border-red-500' : ''}`}>
                            <label class="block text-sm font-medium mb-1">Sensor Type:</label>
                            <select bind:value={editedSensor.type} class="input-field w-full">
                                <option value="" disabled>Select sensor type...</option>
                                {#each Object.values(SensorType) as sensorType}
                                    <option value={sensorType}>{sensorType}</option>
                                {/each}
                            </select>
                        </div>

                        <!-- Refresh Rate Field -->
                        <div class="w-1/4">
                            <label for="refreshRate" class="block text-sm font-medium mb-1"
                                >Refresh Rate:</label
                            >
                            <input
                                type="number"
                                id="refreshRate"
                                bind:value={editedSensor.refresh_rate}
                                class={`input-field w-full ${isInvalidRefreshRate ? 'border-red-500' : ''}`}
                            />
                        </div>
                    </div>

                    <div class="flex justify-end gap-4">
                        <button
                            type="button"
                            class="btn-secondary"
                            on:click={() => {
                                isEditing = false;
                                editedSensor = { ...sensor };
                            }}
                        >
                            Cancel
                        </button>
                        <button type="submit" class="btn-submit">Save</button>
                    </div>
                </form>
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
