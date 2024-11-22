<script lang="ts">
    import { run } from 'svelte/legacy';

    import * as Card from '$lib/components/ui/card';
    import * as Select from '$lib/components/ui/select';
    import NewSensorInput from '$lib/components/FormInput.svelte';
    import { Label } from '@/components/ui/label';
    import type { PageData } from './$types';
    import type { SensorDetails } from '@/types/sensor';
    import { onMount } from 'svelte';
    import { Button } from '@/components/ui/button';
    import { sensorDetailsSchema, sensorTypeSchema } from '$lib/types/sensor';
    import { z } from 'zod';
    import { authFetch } from '@/helpers/fetch';
    import Input from '@/components/ui/input/input.svelte';

    type Props = {
        data: PageData;
        open: boolean;
    };

    let { data, open = $bindable() }: Props = $props();

    let editing = $state(false);
    let loading = $state(true);
    let orgSensor: SensorDetails;
    let sensor: SensorDetails = $state({
        id: '',
        name: '',
        created_at: new Date(),
        type: 'binary_switch',
        active: false,
        refresh_rate: 0,
        uri: '',
    });

    let fieldErrors: Partial<
        Record<'uri' | 'name' | 'refresh_rate' | 'type' | 'active', string>
    > = $state({});

    let globalError: string | null = $state(null);

    const sensorTypes = sensorTypeSchema.options.map((e) => ({
        value: e,
        // TODO: Add capitalization or full lables
        label: e.replace('_', ' '),
    }));

    let selectedType: { value: string; label: string } = $state({
        value: 'not_assigned',
        label: 'not_assigned',
    });

    run(() => {
        if (selectedType) {
            sensor.type = selectedType.value as z.infer<
                typeof sensorTypeSchema
            >;
        }
    });

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

        if (res.ok) {
            console.log(await res.json());
            open = false;
        } else {
            const resJson = await res.json();
            console.log(resJson);
            globalError = resJson.error;
        }
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

        const resJson = await res.json();
        console.log(resJson);

        if (res.ok) {
            open = false;
        } else {
            if (typeof resJson.error === 'string') {
                globalError = resJson.error;
            } else {
                fieldErrors = resJson.error;
            }
        }
    };

    onMount(async () => {
        orgSensor = await data.sensor;
        console.log(orgSensor);
        selectedType =
            sensorTypes.find((e) => e.value === orgSensor.type) ??
            sensorTypes[0];
        sensor = structuredClone(orgSensor);
        loading = false;
    });
</script>

<Card.Root class="w-[512px] border-none shadow-none">
    {#if loading}
        <p>Loading...</p>
    {:else}
        <Card.Header class="text-3xl">
            <Card.Title>Sensor Details</Card.Title>
        </Card.Header>
        <Card.Content class="grid grid-cols-[3fr_4fr] items-center gap-3">
            <NewSensorInput
                name="name"
                label="Name"
                bind:value={sensor.name}
                type="text"
                errors={fieldErrors}
            />
            <NewSensorInput
                name="refresh_rate"
                label="Refresh rate"
                bind:value={sensor.refresh_rate}
                type="number"
                errors={fieldErrors}
            />
            <NewSensorInput
                name="uri"
                label="URI"
                bind:value={sensor.uri}
                type="string"
                errors={fieldErrors}
            />
            <Label
                for="type"
                class="flex items-center justify-between text-base font-semibold"
            >
                Type
                {#if fieldErrors['type']}
                    <span class="text-sm font-normal italic text-red-400"
                        >{fieldErrors['type']}</span
                    >
                {/if}
            </Label>
            <Select.Root bind:selected={selectedType} required name="type">
                <Select.Trigger
                    class={fieldErrors['type'] ? 'border-2 border-red-600' : ''}
                >
                    <Select.Value />
                </Select.Trigger>
                <Select.Content>
                    {#each sensorTypes as type}
                        <Select.Item value={type.value}
                            >{type.label}</Select.Item
                        >
                    {/each}
                </Select.Content>
            </Select.Root>

            <Label
                for="type"
                class="flex items-center justify-between text-base font-semibold"
            >
                Active
                {#if fieldErrors['active']}
                    <span class="text-sm font-normal italic text-red-400"
                        >{fieldErrors['active']}</span
                    >
                {/if}
            </Label>
            <Input
                type="checkbox"
                class="min-w-[4rem] {fieldErrors['active']
                    ? 'border-2 border-red-600'
                    : ''}"
                bind:checked={sensor.active}
            />
        </Card.Content>
        <Card.Footer class="flex justify-end gap-3">
            <div class="flex w-full flex-col items-center justify-center gap-4">
                {#if globalError}
                    <p class="mt-1 text-sm text-red-500">{globalError}</p>
                {/if}

                <div class="flex w-full justify-end gap-3">
                    {#if editing}
                        <Button
                            variant="destructive"
                            size="bold"
                            on:click={handleDelete}>Delete</Button
                        >
                        <Button
                            variant="outline"
                            size="bold"
                            on:click={handleCancel}>Cancel</Button
                        >
                        <Button size="bold" on:click={handleSubmit}
                            >Submit</Button
                        >
                    {:else}
                        <Button
                            on:click={() => {
                                editing = true;
                            }}
                            size="bold">Edit</Button
                        >
                    {/if}
                </div>
            </div>
        </Card.Footer>
    {/if}
</Card.Root>
