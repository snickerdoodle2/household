<script lang="ts">
    import { SensorWebsocket } from '@/helpers/socket.svelte';
    import Chart from './Chart.svelte';
    import Button from '../ui/button/button.svelte';
    import { SvelteMap } from 'svelte/reactivity';

    type Props = {
        sensorId: string;
    };

    let { sensorId }: Props = $props();

    let updateValues = $state(false);

    // websocket stuff
    const ws = new SensorWebsocket();
    $effect(() => {
        if (!ws.ready) return;
        ws.subscribe(sensorId);

        return () => {
            ws.unsubscribe(sensorId);
        };
    });

    let currentSensorOutput = $derived(ws.data.get(sensorId))
    let shownData: SvelteMap<Date, number> | undefined = $state();

    $effect(()=> {
        if (!updateValues) shownData = currentSensorOutput
        console.log([...ws.data.entries()].length)
    })
</script>

<div>
    <div class="p-2">
        {#if shownData}
            <Chart data={shownData} />
        {:else}
            <p>Error opening socket</p>
        {/if}
    </div>

    <div>
        <Button
            on:click={() => {
                updateValues = true;
                shownData = new SvelteMap([...shownData?.entries() ?? []]); //make snapshot
            }}
            size="bold">Stop</Button
        >

        <Button
            on:click={() => {
                updateValues = false;
            }}
            size="bold">Resume</Button
        >

        <Button
            on:click={() => {
                updateValues = false;
                ws.requestSince(sensorId, '24h')
            }}
            size="bold">New</Button
        >
    </div>
</div>
