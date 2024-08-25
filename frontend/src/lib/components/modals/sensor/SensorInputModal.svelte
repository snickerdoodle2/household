<!-- src/routes/AddSensorForm.svelte -->
<script lang="ts">
    import * as Select from '$lib/components/ui/select';
    import { SensorType, type SensorData } from '@/types/Sensor.types';
    import { validateName } from '@/utils/Misc.utils';
    import { closeModal } from '@/utils/Modal.utils';
    import type { Selected } from 'bits-ui';

    export let title = '';
    export let sensorData: SensorData = {
        name: '',
        uri: '',
        type: '',
        refresh_rate: 0,
    };

    export let onSubmit: (data: SensorData) => Promise<void>;

    let isInvalidName: boolean = false;
    let isInvalidType: boolean = false;
    let isInvalidURI: boolean = false;
    let isInvalidRefreshRate: boolean = false;

    const handleSubmit = async (event: Event) => {
        event.preventDefault();
        if (!validateForm()) return;
        await onSubmit(sensorData);
        closeModal();
    };

    function validateForm(): boolean {
        let isInvalid = false;
        const uriRegex =
            /^(25[0-5]|2[0-4]\d|1\d\d|\d\d?)\.(25[0-5]|2[0-4]\d|1\d\d|\d\d?)\.(25[0-5]|2[0-4]\d|1\d\d|\d\d?)\.(25[0-5]|2[0-4]\d|1\d\d|\d\d?):([1-9]\d{0,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d|655[0-2]\d|6553[0-5])$/;

        if (validateName(sensorData.name).isError) {
            isInvalidName = true;
            isInvalid = true;
        } else {
            isInvalidName = false;
        }

        if (sensorData.type === '') {
            isInvalidType = true;
            isInvalid = true;
        } else {
            isInvalidType = false;
        }

        if (!uriRegex.test(sensorData.uri)) {
            isInvalidURI = true;
            isInvalid = true;
        } else {
            isInvalidURI = false;
        }

        if (sensorData.refresh_rate <= 0) {
            isInvalidRefreshRate = true;
            isInvalid = true;
        } else {
            isInvalidRefreshRate = false;
        }
        return !isInvalid;
    }

    const updateSelected = (item: Selected<string> | undefined) => {
        if (!item || item.value.length === 0) {
            return;
        }
        sensorData.type = item.value;
    };
</script>

<main>
    <div class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
        <div class="bg-background rounded-lg shadow-lg w-full max-w-lg p-6 relative">
            <button type="button" class="absolute top-2 right-2 btn-exit" on:click={closeModal}>
                &times;
            </button>
            <h2 class="text-2xl font-bold mb-4">{title}</h2>
            <form on:submit={handleSubmit} class="space-y-4">
                <div>
                    <label for="name" class="block text-sm font-medium mb-1">Name:</label>
                    <input
                        type="text"
                        id="name"
                        bind:value={sensorData.name}
                        class={`input-field w-full ${isInvalidName ? 'border-red-500' : ''}`}
                    />
                </div>
                <div>
                    <label for="uri" class="block text-sm font-medium mb-1">URI:</label>
                    <input
                        type="text"
                        id="uri"
                        bind:value={sensorData.uri}
                        class={`input-field w-full ${isInvalidURI ? 'border-red-500' : ''}`}
                    />
                </div>
                <div class="flex gap-4">
                    <div class={`w-3/4 ${isInvalidType ? 'border-red-500' : ''}`}>
                        <label class="block text-sm font-medium mb-1">Sensor Type:</label>
                        <select bind:value={sensorData.type} class="input-field w-full">
                            <option value="" disabled>Select sensor type...</option>
                            {#each Object.values(SensorType) as sensorType}
                                <option value={sensorType}>{sensorType}</option>
                            {/each}
                        </select>
                    </div>

                    <div class="w-1/4">
                        <label for="refreshRate" class="block text-sm font-medium mb-1"
                            >Refresh Rate:</label
                        >
                        <input
                            type="number"
                            id="refreshRate"
                            bind:value={sensorData.refresh_rate}
                            class={`input-field w-full ${isInvalidRefreshRate ? 'border-red-500' : ''}`}
                        />
                    </div>
                </div>
                <button type="submit" class="btn-submit w-full">Submit</button>
            </form>
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
