<script lang="ts">
    import * as Card from '$lib/components/ui/card';
    import { Input } from '@/components/ui/input';
    import { Label } from '@/components/ui/label';
    import type { PageData } from './$types';
    import type { SensorDetails } from '@/types/sensor';
    import { onMount } from 'svelte';
    import { Button } from '@/components/ui/button';
    export let data: PageData;

    let editing = false;
    let loading = true;
    let orgSensor: SensorDetails;
    let sensor: SensorDetails;

    onMount(async () => {
        orgSensor = await data.sensor;
        sensor = orgSensor;
        loading = false;
    });
</script>

<Card.Root class="w-[512px] border-none">
    {#if loading}
        <p>Loading...</p>
    {:else}
        <Card.Header class="text-3xl">
            <Card.Title>Sensor Details</Card.Title>
        </Card.Header>
        <Card.Content
            class="grid grid-cols-2 items-center gap-3 text-xl font-semibold"
        >
            <Label for="name">Name</Label>
            <Input
                type="text"
                name="name"
                disabled={!editing}
                bind:value={sensor.name}
            />
            <Label for="refresh_rate">Refresh Rate</Label>
            <Input
                type="number"
                name="refresh_rate"
                disabled={!editing}
                bind:value={sensor.refresh_rate}
            />
            <Label for="uri">URI</Label>
            <Input
                type="text"
                name="uri"
                disabled={!editing}
                bind:value={sensor.uri}
            />
            <Label for="sensor_type">Type</Label>
        </Card.Content>
        <Card.Footer class="flex justify-end">
            {#if editing}
                <Button variant="destructive">Delete</Button>
                <Button variant="outline">Cancel</Button>
                <Button>Submit</Button>
            {:else}
                <Button
                    on:click={() => {
                        editing = true;
                    }}>Edit</Button
                >
            {/if}
        </Card.Footer>
    {/if}
</Card.Root>
