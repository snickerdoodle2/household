<script lang="ts">
    import { type Sensor } from '@/types/sensor';
    import { DotsVertical } from 'svelte-radix';
    import { SensorWebsocket } from '@/helpers/socket.svelte';
    type Props = {
        sensor: Sensor;
    };

    let { sensor }: Props = $props();
    const ws = new SensorWebsocket();

    $effect(() => {
        if (!ws.ready) return;
        ws.subscribe(sensor.id);

        return () => {
            ws.unsubscribe(sensor.name);
        };
    });

    let data = $derived(ws.data.get(sensor.id));
</script>

<div class="flex flex-col gap-2 rounded-lg bg-accent px-4 py-2">
    <div class="flex items-center justify-between">
        <span class="text-xl">{sensor.name} </span>
        <div class="flex items-center gap-2">
            <div class={`aspect-square w-2 rounded-full`}></div>
            <a href={`/details/${sensor.id}`}
                ><DotsVertical class="h-5 w-5" /></a
            >
        </div>
    </div>
    {#if data}
        {(
            data.values().reduce((acc, cur) => acc + cur, 0) /
            data.values().reduce((acc) => acc + 1, 0)
        ).toFixed(2)} ({data
            .values()
            .reduce((acc, cur) => acc + cur, 0)
            .toFixed()})
    {:else}
        <p>Error opening socket</p>
    {/if}
</div>
