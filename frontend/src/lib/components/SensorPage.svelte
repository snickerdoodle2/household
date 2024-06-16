<script lang="ts">
    import { SensorType, type Sensor } from '@/types/sensor';
    import { onMount } from 'svelte';
    import {
        Table,
        TableHead,
        TableHeadCell,
        TableBody,
        TableBodyRow,
        TableBodyCell,
        Spinner,
        Button,
        Dropdown,
        Checkbox,
    } from 'flowbite-svelte';
    import { sensors } from '@/stores/stores';

    let loading = true;
    let error: string | null = null;

    async function fetchSensors() {
        try {
            const response = await fetch('http://localhost:8080/api/v1/sensor');
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();

            sensors.set(data.data);
            console.log('Fetched sensors:', sensors);
        } catch (err) {
            error = 'Failed to fetch sensors';
            console.error('Failed to fetch sensors:', err);
        } finally {
            loading = false;
        }
    }

    // Function to delete a sensor
    async function deleteSensor(id: string) {
        console.log('Deleting sensor with id:', id);
        try {
            const response = await fetch(
                `http://localhost:8080/api/v1/sensor/${id}`,
                {
                    method: 'DELETE',
                }
            );

            if (response.ok) {
            } else {
                console.error('Error deleting sensor:', response.statusText);
            }
        } catch (error) {
            console.error('Error deleting sensor:', error);
        }

        fetchSensors();
    }

    onMount(fetchSensors);

    const sensorVisibility = {
        [SensorType.BINARY_SENSOR]: true,
        [SensorType.BINARY_SWITCH]: true,
        [SensorType.BUTTON]: true,
        [SensorType.DECIMAL_SENSOR]: true,
        [SensorType.DECIMAL_SWITCH]: true,
    };
</script>

<main>
    {#if loading}
        <Spinner />
    {:else if error}
        <p class="error">{error}</p>
    {:else}
        <div class="p-2 place-content-end">
            <script>
                import {
                    Button,
                    Dropdown,
                    DropdownItem,
                    Checkbox,
                } from 'flowbite-svelte';
                import { ChevronDownOutline } from 'flowbite-svelte-icons';
            </script>

            <Button
                class="bg-orange-500 hover:bg-orange-700 text-white font-bold rounded"
                >Filter</Button
            >

            <Dropdown placement={"right"} class="w-48 p-3 space-y-1 text-sm">
                {#each Object.values(SensorType) as type}
                    <li
                        class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600"
                    >
                        <Checkbox
                            id={type}
                            bind:checked={sensorVisibility[type]}
                        >
                            {type}
                        </Checkbox>
                    </li>
                {/each}
            </Dropdown>
        </div>

        <div class="table-container">
            <Table class="tables" hoverable={true}>
                <TableHead>
                    <TableHeadCell>ID</TableHeadCell>
                    <TableHeadCell>Name</TableHeadCell>
                    <TableHeadCell>URI</TableHeadCell>
                    <TableHeadCell>Type</TableHeadCell>
                    <TableHeadCell>Refresh Rate</TableHeadCell>
                    <TableHeadCell>Created At</TableHeadCell>
                    <TableHeadCell>Version</TableHeadCell>
                    <TableHeadCell>Actions</TableHeadCell>
                </TableHead>
                <TableBody tableBodyClass="divide-y">
                    {#each $sensors as sensor}
                        {#if sensorVisibility[sensor.type]}
                            <TableBodyRow>
                                <TableBodyCell>{sensor.id}</TableBodyCell>
                                <TableBodyCell>{sensor.name}</TableBodyCell>
                                <TableBodyCell>{sensor.uri}</TableBodyCell>
                                <TableBodyCell>{sensor.type}</TableBodyCell>
                                <TableBodyCell
                                    >{sensor.refresh_rate}</TableBodyCell
                                >
                                <TableBodyCell
                                    >{sensor.created_at}</TableBodyCell
                                >
                                <TableBodyCell>{sensor.version}</TableBodyCell>
                                <TableBodyCell>
                                    <Button
                                        on:click={() => console.log('Edit')}
                                        color="blue"
                                        class="mr-2">Edit</Button
                                    >
                                    <Button
                                        on:click={() => deleteSensor(sensor.id)}
                                        color="red"
                                        class="mr-2">Remove</Button
                                    >
                                    <Button
                                        on:click={() => console.log('Monitor')}
                                        color="green">Monitor</Button
                                    >
                                </TableBodyCell>
                            </TableBodyRow>
                        {/if}
                    {/each}
                </TableBody>
            </Table>
        </div>
    {/if}
</main>

<style>
    main {
        display: flex;
        flex-direction: column;
        height: 100%;
        width: 100%;
    }

    .table-container {
        min-width: 100%;
        flex: 1;
        overflow: auto; /* Enable scrolling if content overflows */
        max-width: 100%; /* Ensure the table doesn't exceed the container's width */
        height: 100%;
        top: 0;
        left: 0;
        margin: 0 auto; /* Center the container */
    }

    .error {
        color: red;
    }
</style>
