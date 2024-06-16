<!-- src/routes/AddSensorForm.svelte -->
<script lang="ts">
    import type { SensorData } from '@/types/sensor';
    import SensorInputModal from './SensorInputModal.svelte';

    export let isOpen = true;

    export let afterSubmit: () => Promise<void> = async () => {};

    async function addSensor({ name, uri, type, refresh_rate }: SensorData) {
        try {
            const response = await fetch(
                'http://localhost:8080/api/v1/sensor',
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
        afterSubmit();
    }
</script>

<main>
    <SensorInputModal title={"Add New Device"} onSubmit={addSensor} bind:isOpen />
</main>

<style>
</style>
