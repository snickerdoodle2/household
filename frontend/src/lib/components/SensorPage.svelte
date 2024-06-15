<script lang="ts">
    import { SensorType, type Sensor } from '@/types/sensor';
    import { onMount } from 'svelte';
    import { Checkbox, Table, TableHead, TableHeadCell, TableBody, TableBodyRow, TableBodyCell, Spinner, Button } from 'flowbite-svelte';
  
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
        <Spinner />
    {:else if error}
      <p class="error">{error}</p>
    {:else}
      <div class="filter">
        <Table  bordered={true} class="w-full text-sm text-left text-gray-500">
            <TableHead> 
                <TableHeadCell>Filter by type</TableHeadCell>
            </TableHead>
            <TableBody>
                {#each Object.values(SensorType) as type}
                    <TableBodyRow>
                        <Checkbox id={type} bind:checked={selectedTypes[type]}>{type}</Checkbox>
                    </TableBodyRow>
                {/each}
            </TableBody>
          </Table>

          
      </div>
  
      <div class="table-container">
        <Table hoverable={true} >
          <TableHead>
            <TableHeadCell></TableHeadCell>
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
            {#each filteredSensors() as sensor}
              {#if selectedTypes[sensor.type]}
                <TableBodyRow>
                  <TableBodyCell></TableBodyCell>
                  <TableBodyCell>{sensor.id}</TableBodyCell>
                  <TableBodyCell>{sensor.name}</TableBodyCell>
                  <TableBodyCell>{sensor.uri}</TableBodyCell>
                  <TableBodyCell>{sensor.type}</TableBodyCell>
                  <TableBodyCell>{sensor.refresh_rate}</TableBodyCell>
                  <TableBodyCell>{sensor.created_at}</TableBodyCell>
                  <TableBodyCell>{sensor.version}</TableBodyCell>
                  <TableBodyCell>
                    <Button on:click={() => console.log("Edit")} color="blue" class="mr-2">Edit</Button>
                    <Button on:click={() => console.log("Remove")} color="red" class="mr-2">Remove</Button>
                    <Button on:click={() => console.log("Monitor")} color="green">Monitor</Button>
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
    }

    .filter {
        top: 0;
        left: 0;
    }
    .table-container {
        flex: 1;
        overflow: auto; /* Enable scrolling if content overflows */
    }
  
    .error {
        color: red;
    }
  </style>
  