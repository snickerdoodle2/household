<!-- src/routes/AddSensorForm.svelte -->
<script lang="ts">
    import type { SensorData } from '@/types/Sensor.types';
    import SensorInputModal from './SensorInputModal.svelte';
    import { openedModalStore } from '@/stores/Stores';
    import { ModalType, isModalData } from '@/types/Modal.types';
    import { SERVER_URL } from '@/config/const';

    let open = false;
    let id: string = '';

    openedModalStore.subscribe((value) => {
        if (isModalData(ModalType.MODIFY_SENSOR, value)) {
            open = true;
            id = value.data.id;
        } else {
            open = false;
        }
    });

    export let onClose: () => Promise<void> = async () => {};

    export const sensorData: SensorData = {
        name: '',
        uri: '',
        type: '',
        refresh_rate: 0,
    };

    async function modifySensor({
        name,
        uri,
        type,
        refresh_rate,
    }: SensorData) {
        try {
            const response = await fetch(`${SERVER_URL}/api/v1/sensor/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ name, uri, type, refresh_rate }),
            });

            console.log(`${SERVER_URL}/api/v1/sensor/${id}`);
            if (!response.ok) {
                // TODO: error json handling
                const errorData = await response.json();
                console.error('Error:', errorData);
                alert(`Error: ${errorData.error}`);
            } else {
                // TODO: nice pop-up window instead of alert
                const responseData = await response.json();
                console.log('Success:', responseData);
                alert('Sensor modified successfully!');
            }
        } catch (error) {
            console.error('Network Error:', error);
            alert('Network error. Please try again later.');
        }
        onClose();
    }
</script>

<main>
    {#if isModalData(ModalType.MODIFY_SENSOR, $openedModalStore)}
        <SensorInputModal
            title={'Modify Device'}
            onSubmit={modifySensor}
            bind:sensorData={$openedModalStore.data}
            bind:open
        />
    {/if}
</main>

<style>
</style>
