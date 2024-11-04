<script lang="ts">
    import { Chart } from 'chart.js/auto';
    import { onDestroy, onMount } from 'svelte';
    import { type SocketStore } from '$lib/helpers/socket';
    import { get } from 'svelte/store';
    type Props = {
        socket: SocketStore;
    };

    let { socket }: Props = $props();
    let chartEl: HTMLCanvasElement = $state();
    let chart: Chart;

    const unsubscribe = socket.subscribe((data) => {
        if (!chart) return;
        if (!data) return;
        if (!chart.data.labels) return;
        chart.data.labels.push((chart.data.labels.at(-1) as number) + 1);
        chart.data.datasets[0].data.push(data.values.at(-1) as number);

        chart.data.labels.shift();
        chart.data.datasets[0].data.shift();
        chart.update();
    });

    onDestroy(() => {
        unsubscribe();
    });

    onMount(() => {
        if (!$socket) return;
        const data = get(socket);
        if (!data) return;

        chart = new Chart(chartEl, {
            type: 'line',
            data: {
                labels: data.values.map((_, i) => i),
                datasets: [
                    {
                        label: 'lol',
                        data: data.values,
                    },
                ],
            },
        });
    });
</script>

<canvas bind:this={chartEl}></canvas>
