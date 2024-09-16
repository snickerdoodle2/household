<script lang="ts">
    import * as Card from '$lib/components/ui/card';
    import * as Select from '$lib/components/ui/select';
    import { Input } from '@/components/ui/input';
    import { Label } from '@/components/ui/label';
    import type { PageData } from './$types';
    import type { SensorDetails } from '@/types/sensor';
    import { onMount } from 'svelte';
    import { Button } from '@/components/ui/button';
    import { sensorDetailsSchema, sensorTypeSchema } from '$lib/types/sensor';
    import { z } from 'zod';
    import { authFetch } from '@/helpers/fetch';

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

    const labelClass = 'font-semibold text-base';

    const handleCancel = () => {
        sensor = structuredClone(orgSensor);
        selectedType =
            sensorTypes.find((e) => e.value === orgSensor.type) ??
            sensorTypes[0];
        editing = false;
    };

    const handleDelete = async () => {
        // TODO: ask for confirmation!!!
        const res = await authFetch(`/api/v1/sensor/${orgSensor.id}`, {
            method: 'DELETE',
        });
        console.log(await res.json());
    };

    const handleSubmit = async () => {
        const { data, success, error } = sensorDetailsSchema.safeParse({
            ...sensor,
            refresh_rate: +sensor.refresh_rate,
        });
        if (!success) {
            console.log(error.issues);
            return;
        }

        const { id, created_at, ...rest } = data; // eslint-disable-line @typescript-eslint/no-unused-vars

        const res = await authFetch(`/api/v1/sensor/${orgSensor.id}`, {
            method: 'PUT',
            body: JSON.stringify(rest),
        });

        console.log(await res.json());
    };

    onMount(async () => {
        orgSensor = await data.sensor;
        selectedType =
            sensorTypes.find((e) => e.value === orgSensor.type) ??
            sensorTypes[0];
        sensor = structuredClone(orgSensor);
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
        <Card.Content class="grid grid-cols-[1fr_2fr] items-center gap-3">
            <Label for="name" class={labelClass}>Name</Label>
            <Input
                type="text"
                name="name"
                disabled={!editing}
                bind:value={sensor.name}
            />
            <Label for="refresh_rate" class={labelClass}>Refresh Rate</Label>
            <Input
                type="number"
                name="refresh_rate"
                disabled={!editing}
                bind:value={sensor.refresh_rate}
            />
            <Label for="uri" class={labelClass}>URI</Label>
            <Input
                type="text"
                name="uri"
                disabled={!editing}
                bind:value={sensor.uri}
            />
            <Label for="sensor_type" class={labelClass}>Type</Label>
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
        <Card.Footer class="flex justify-end gap-3">
            {#if editing}
                <Button
                    variant="destructive"
                    size="bold"
                    on:click={handleDelete}>Delete</Button
                >
                <Button variant="outline" size="bold" on:click={handleCancel}
                    >Cancel</Button
                >
                <Button size="bold" on:click={handleSubmit}>Submit</Button>
            {:else}
                <Button
                    on:click={() => {
                        editing = true;
                    }}
                    size="bold">Edit</Button
                >
            {/if}
        </Card.Footer>
    {/if}
</Card.Root>
