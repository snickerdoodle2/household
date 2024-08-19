<!-- src/routes/AddSensorForm.svelte -->
<script lang="ts">
    import type { SensorData } from '@/types/Sensor.types';
    import SensorInputModal from './SensorInputModal.svelte';
    import { ModalType } from '@/types/Modal.types';
    import { openedModalStore } from '@/stores/Stores';
    import { SERVER_URL } from '@/config/const';
    import { isModalData } from '@/utils/Modal.utils';

    export let onClose: () => Promise<void> = async () => {};

    async function addSensor({ name, uri, type, refresh_rate }: SensorData) {
        try {
            const response = await fetch(
                `${SERVER_URL}/api/v1/sensor`,
                {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        name,
                        uri,
                        type,
                        refresh_rate,
                    }),
                }
            );

            if (!response.ok) {
                // TODO: error json handling
                const errorData = await response.json();
                console.error('Error:', errorData);
                alert(`Error: ${errorData.error}`);
            } else {
                // TODO: nice pop-up window instead of alert
                const responseData = await response.json();
                console.log('Success:', responseData);
                alert('Sensor added successfully!');
            }
        } catch (error) {
            console.error('Network Error:', error);
            alert('Network error. Please try again later.');
        }
        onClose();
    }
</script>

<main>
    {#if isModalData(ModalType.ADD_SENSOR, $openedModalStore)}
        <SensorInputModal title={'Add New Device'} onSubmit={addSensor}/>
    {/if}
</main>

<style>
</style>
