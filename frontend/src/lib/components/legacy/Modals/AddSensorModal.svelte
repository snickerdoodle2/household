<!-- src/routes/AddSensorForm.svelte -->
<script lang="ts">
    import type { SensorData } from '@/types/sensor';
    import SensorInputModal from './SensorInputModal.svelte';
    import { ModalType, isModalData } from '@/types/modal';
    import { openedModalStore } from '@/stores/stores';
    import { SERVER_URL } from '@/const';

    export let onClose: () => Promise<void> = async () => {};

    let open = false;

    openedModalStore.subscribe((value) => {
        open = isModalData(ModalType.ADD_SENSOR, value);
    });

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
        <SensorInputModal title={'Add New Device'} onSubmit={addSensor} bind:open />
    {/if}
</main>

<style>
</style>
