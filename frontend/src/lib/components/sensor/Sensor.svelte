<script lang="ts">
    import { createBubbler, preventDefault } from 'svelte/legacy';

    const bubble = createBubbler();
    import { type Sensor } from '@/types/sensor';
    import { DotsVertical } from 'svelte-radix';
    import { socketStore } from '$lib/helpers/socket';
    import { onDestroy } from 'svelte';
    import Chart from './Chart.svelte';
    interface Props {
        sensor: Sensor;
    }

    let { sensor }: Props = $props();

    let socket = socketStore(sensor.id);

    onDestroy(() => {
        socket.close();
    });
</script>

<div class="flex flex-col gap-2 rounded-lg bg-accent px-4 py-2">
    {#if $socket}
        <div class="flex items-center justify-between">
            <span class="text-xl">{sensor.name} </span>
            <div class="flex items-center gap-2">
                <div
                    class={`aspect-square w-2 rounded-full ${$socket.status === 'ONLINE' ? 'bg-green-400' : 'bg-red-400'}`}
                ></div>
                <a
                    href={`/details/${sensor.id}`}
                    onclick={preventDefault(bubble('click'))}
                    ><DotsVertical class="h-5 w-5" /></a
                >
            </div>
        </div>
        <Chart {socket} />
    {:else}
        <p>Error opening socket</p>
    {/if}
</div>
