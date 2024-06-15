<script lang="ts">
  import { SensorType, type Sensor } from '@/types/sensor';
    import { onMount } from 'svelte';

    let sensors: Sensor[] = [];
    let loading = true;
    let error: string | null = null;

    async function fetchSensors() {
        try {
            const response = await fetch('http://localhost:8080/api/v1/sensor');
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();
            sensors = data.data;
            console.log('Fetched sensors:', sensors);
        } catch (err) {
            error = 'Failed to fetch sensors';
            console.error('Failed to fetch sensors:', err);
        } finally {
            loading = false;
        }
    }

    onMount(fetchSensors);

    const selectedTypes: Record<SensorType, boolean> = {
        [SensorType.BINARY_SENSOR]: true,
        [SensorType.BINARY_SWITCH]: true,
        [SensorType.BUTTON]: true,
        [SensorType.DECIMAL_SENSOR]: true,
        [SensorType.DECIMAL_SWITCH]: true,
    }

    function filteredSensors() {
        return sensors.filter(sensor => selectedTypes[sensor.type]);
    }
</script>


<main>
    {#if loading}
        <p>Loading...</p>
    {:else if error}
        <p class="error">{error}</p>
    {:else}
        <div class="filter">
            <h2>Filter Sensors</h2>
            {#each Object.values(SensorType) as type}
                <div>
                    <label>
                        <input type="checkbox" bind:checked={selectedTypes[type]} />
                        {type}
                    </label>
                </div>
            {/each}
        </div>

        <table>
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>URI</th>
                    <th>Type</th>
                    <th>Refresh Rate</th>
                    <th>Created At</th>
                    <th>Version</th>
                </tr>
            </thead>
            <tbody>
                {#each filteredSensors() as sensor}
                    {#if selectedTypes[sensor.type]}
                        <tr>
                            <td>{sensor.id}</td>
                            <td>{sensor.name}</td>
                            <td>{sensor.uri}</td>
                            <td>{sensor.type}</td>
                            <td>{sensor.refresh_rate}</td>
                            <td>{sensor.created_at}</td>
                            <td>{sensor.version}</td>
                        </tr>
                    {/if}
                {/each}
            </tbody>
        </table>
    {/if}
</main>

<style>
    main {
        padding: 2rem;
    }
    table {
        width: 100%;
        border-collapse: collapse;
        margin-top: 1rem;
    }
    th, td {
        border: 1px solid #ddd;
        padding: 8px;
    }
    th {
        background-color: #616161;
        color: white;
    }
    .error {
        color: red;
    }
    .filter {
        margin-bottom: 1rem;
    }
</style>