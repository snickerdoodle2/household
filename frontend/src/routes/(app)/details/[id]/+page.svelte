<script lang="ts">
    import * as Card from '$lib/components/ui/card';
    import * as Select from '$lib/components/ui/select';
    import { Input } from '@/components/ui/input';
    import { Label } from '@/components/ui/label';
    import type { PageData } from './$types';
    import type { SensorDetails } from '@/types/sensor';
    import { onMount } from 'svelte';
    import { Button } from '@/components/ui/button';
    import { sensorTypeSchema } from '$lib/types/sensor';
    import { z } from 'zod';

    export let data: PageData;

    let editing = false;
    let loading = true;
    let orgSensor: SensorDetails;
    let sensor: SensorDetails;

    const sensorTypes = sensorTypeSchema.options.map((e) => ({
        value: e,
        // TODO: Add capitalization or full lables
        label: e.replace('_', ' '),
    }));

    let selectedType: { value: string; label: string };

    $: {
        if (selectedType) {
            sensor.type = selectedType.value as z.infer<
                typeof sensorTypeSchema
            >;
        }
    }

    onMount(async () => {
        orgSensor = await data.sensor;
        selectedType =
            sensorTypes.find((e) => e.value === orgSensor.type) ??
            sensorTypes[0];
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
            <Select.Root disabled={!editing} bind:selected={selectedType}>
                <Select.Trigger>
                    <Select.Value />
                </Select.Trigger>
                <Select.Content>
                    {#each sensorTypes as type}
                        <Select.Item value={type.label}
                            >{type.label}</Select.Item
                        >
                    {/each}
                </Select.Content>
            </Select.Root>
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
