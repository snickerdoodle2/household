<script lang="ts">
    // TODO: calculate grid size
    import type { LayoutData } from './$types';
    import { goto } from '$app/navigation';
    import { Plus } from 'svelte-radix';
    import { Button } from '@/components/ui/button';
    import * as Select from '$lib/components/ui/select';
    import Label from '@/components/ui/label/label.svelte';
    import { EyeOpen } from 'svelte-radix';
    import type { Sensor as SensorType } from '@/types/sensor';
    import Sensor from '@/components/sensor/Sensor.svelte';
    import { authFetch } from '@/helpers/fetch';
    import { onMount } from 'svelte';

    const handleCreate = () => {
        goto(window.location.href + 'create');
    };

    type Props = {
        data: LayoutData;
    };

    let { data }: Props = $props();

    let sensors: SensorType[] | undefined = $state(undefined);

    const showSensor = async (sensor: SensorType) => {
        sensor.hidden = false;
        const res = await authFetch(`/api/v1/sensor/${sensor.id}/hidden`, {
            method: 'PUT',
            body: JSON.stringify({
                hidden: false
            }),
        });

        const resJson = await res.json();
        console.log(resJson);

        if (!res.ok) {
            console.log(resJson.error);
        }
    };

    const showAll = async () => {
        if (!sensors) return;
        for (const sensor of sensors) {
            showSensor(sensor);
        }
    };

    onMount(async () => {
        sensors = await data.sensors;
    });
</script>

{#if sensors}
    <div class="absolute right-10 top-10">
        <div class="flex gap-4">
            <Label class="text-xl">Show sensors</Label>
            <Select.Root required name="type">
                <Select.Trigger class="min-w-64"></Select.Trigger>
                <Select.Content class="min-w-32">
                    {#each sensors as sensor}
                        {#if sensor.hidden}
                            <Select.Item
                                showSelected={false}
                                value={sensor.id}
                                on:click={() => showSensor(sensor)}
                            >
                                <div
                                    class="w-full flex items-center justify-between"
                                >
                                    <Label class="text-md">{sensor.name}</Label>
                                    <span>
                                        <EyeOpen
                                            on:click={() => showSensor(sensor)}
                                        />
                                    </span>
                                </div>
                            </Select.Item>
                        {/if}
                    {/each}
                </Select.Content>
            </Select.Root>
            <Button on:click={showAll}>Show all</Button>
        </div>
    </div>
    <div class="flex h-full items-start gap-4 md:py-20">
        <div
            class="grid flex-1 grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5"
        >
            {#each sensors as sensor, i}
                {#if !sensor.hidden}
                    <Sensor bind:sensor={sensors[i]} />
                {/if}
            {/each}
            <div class="flex items-center justify-center">
                <Button variant="outline" size="icon" on:click={handleCreate}>
                    <Plus />
                </Button>
            </div>
        </div>
    </div>
{/if}
