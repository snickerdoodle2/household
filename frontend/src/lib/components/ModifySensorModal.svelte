<!-- src/routes/AddSensorForm.svelte -->
<script lang="ts">
    import type { SensorData } from '@/types/sensor';
    import SensorInputModal from './SensorInputModal.svelte';
    import { ModifySensorModalData } from '@/stores/stores';

    export let isOpen = true;

    export let id: string = '';

    export let afterSubmit: () => Promise<void> = async () => {};

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
        refresh_rate: refreshRate,
    }: SensorData) {
        try {
            const response = await fetch(
                `http://localhost:8080/api/v1/sensor/${id}`,
                {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ name, uri, type, refreshRate }),
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
                alert('Sensor modified successfully!');
            }
        } catch (error) {
            console.error('Network Error:', error);
            alert('Network error. Please try again later.');
        }
        afterSubmit();
    }
</script>

<main>
    <SensorInputModal onSubmit={modifySensor} bind:isOpen bind:sensorData={$ModifySensorModalData}/>
</main>

<style>
</style>
