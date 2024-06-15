<script lang="ts">
    import { onMount } from 'svelte';

    interface Sensor {
        id: string;
        name: string;
        uri: string;
        type: string;
        refresh_rate: number;
        created_at: string;
        version: number;
    }

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
</script>


<main>
    {#if loading}
        <p>Loading...</p>
    {:else if error}
        <p class="error">{error}</p>
    {:else}
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
                {#each sensors as sensor}
                    <tr>
                        <td>{sensor.id}</td>
                        <td>{sensor.name}</td>
                        <td>{sensor.uri}</td>
                        <td>{sensor.type}</td>
                        <td>{sensor.refresh_rate}</td>
                        <td>{sensor.created_at}</td>
                        <td>{sensor.version}</td>
                    </tr>
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
    }
    .error {
        color: red;
    }
</style>