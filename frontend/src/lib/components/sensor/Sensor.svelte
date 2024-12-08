<script lang="ts">
    import { type Sensor } from '@/types/sensor';
    import { DotsVertical, EyeNone } from 'svelte-radix';
    import { AppWebsocket } from '@/helpers/socket.svelte';
    import Chart from './Chart.svelte';
    import Button from '../ui/button/button.svelte';
    import { authFetch } from '@/helpers/fetch';
    type Props = {
        sensor: Sensor;
    };

    let { sensor = $bindable() }: Props = $props();
    const ws = new AppWebsocket();

    $effect(() => {
        if (!ws.ready) return;
        ws.subscribe(sensor.id);

        return () => {
            ws.unsubscribe(sensor.name);
        };
    });

    const hideSensor = async () => {
        sensor.hidden = true;
        const res = await authFetch(`/api/v1/sensor/${sensor.id}/hidden`, {
            method: 'PUT',
            body: JSON.stringify({
                hidden: true,
            }),
        });

        const resJson = await res.json();
        console.log(resJson);

        if (!res.ok) {
            console.log(resJson.error);
        }
    };

    let data = $derived(ws.data.get(sensor.id));
</script>

<div class="flex flex-col gap-2 rounded-lg bg-accent px-4 py-2">
    <div class="flex items-center justify-between">
        <span class="text-xl">{sensor.name} </span>
        <div class="flex items-center gap-2">
            <div class={`aspect-square w-2 rounded-full`}></div>
            <!-- TODO: bubble up on:click to show modal -->
            <Button variant="ghost" onclick={hideSensor}>
                <EyeNone class="h-5 w-5" />
            </Button>
            <a href={`/details/${sensor.id}`}
                ><DotsVertical class="h-5 w-5" /></a
            >
        </div>
    </div>
    {#if data}
        <a href={`/view/${sensor.id}`}>
            <Chart {data} fixedView={null} />
        </a>
    {:else}
        <p>Error opening socket</p>
    {/if}
</div>
