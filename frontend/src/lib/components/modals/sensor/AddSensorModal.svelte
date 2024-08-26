<!-- src/routes/AddSensorForm.svelte -->
<script lang="ts">
    import type { SensorData } from '@/types/Sensor.types';
    import SensorInputModal from './SensorInputModal.svelte';
    import { addSensor } from '@/utils/requests/Sensor.requests';
    import { syncSensorConfig } from '@/utils/Sync.utils';

    async function handleSensorAddition(data: SensorData) {
        const result = await addSensor(data);
        if (result.isError) {
            console.log(`Failed to add a sensor: ${result.error}`);
        } else {
            console.log('Sensor added successfully', result.data);
            syncSensorConfig();
        }
    }
</script>

<main>
    <SensorInputModal title={'Add New Device'} onSubmit={handleSensorAddition} />
</main>

<style>
</style>
