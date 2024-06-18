<!-- src/routes/AddSensorForm.svelte -->
<script lang="ts">
    import * as Select from '$lib/components/ui/select';
    import { closeModal } from '@/stores/stores.utils';
    import { SensorType, type SensorData } from '@/types/sensor';
    import type { Selected } from 'bits-ui';
    import { Modal } from 'flowbite-svelte';

    export let title = '';
    export let open = false;
    export let sensorData: SensorData = {
        name: '',
        uri: '',
        type: '',
        refresh_rate: 0,
    };

    export let onSubmit: (data: SensorData) => Promise<void>;

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

<Modal {title} autoclose={false} bind:open>
    <form on:submit={handleSubmit}>
        <div>
            <label for="name">Name:</label>
            <input type="text" id="name" bind:value={sensorData.name} />
        </div>
        <div>
            <label for="uri">URI:</label>
            <input
                type="text"
                id="uri"
                bind:value={sensorData.uri}
                class={isInvalidURI ? 'invalid' : ''}
            />
        </div>
        <div class="row">
            <div class="select-container {isInvalidType ? 'invalid' : ''}">
                <Select.Root onSelectedChange={updateSelected}>
                    <Select.Trigger class="max-w-96">
                        <Select.Value placeholder="Select sensor type..." />
                    </Select.Trigger>
                    <Select.Content>
                        {#each Object.values(SensorType) as sensor}
                            <Select.Item value={sensor} label={sensor}>
                                {sensor}
                            </Select.Item>
                        {/each}
                    </Select.Content>
                </Select.Root>
            </div>
            <div class="refresh-rate-container">
                <label for="refreshRate">Refresh Rate:</label>
                <input
                    type="number"
                    id="refreshRate"
                    bind:value={sensorData.refresh_rate}
                    class={isInvalidRefreshRate ? 'invalid' : ''}
                />
            </div>
        </div>
        <button type="submit">Submit</button>
    </form>
</Modal>

<style>
    form {
        display: flex;
        flex-direction: column;
        gap: 1em;
    }
    .row {
        display: flex;
        gap: 1em;
        align-items: flex-end;
        flex-flow: nowrap;
    }
    .select-container {
        width: 75%;
    }
    .refresh-rate-container {
        width: 25%;
    }
    .invalid {
        border-color: red;
        border-style: solid;
        border-width: 2px;
        border-radius: 5px;
    }
    div {
        display: flex;
        flex-direction: column;
    }
    label {
        margin-bottom: 0.5em;
    }
    input {
        padding: 0.5em;
        font-size: 1em;
    }
    button {
        align-self: center;
        border-style: solid;
        border-width: 2px;
        border-color: white;
        border-radius: 5px;
        padding: 0.5em 1em;
        font-size: 1em;
        cursor: pointer;
    }
</style>
