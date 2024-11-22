<script lang="ts">
import { Chart } from 'chart.js/auto';
import { onMount } from 'svelte';
import type { SvelteMap } from 'svelte/reactivity';
const MAX_RECORDS = 32;

let ctx: HTMLCanvasElement;
let chart: Chart;
let mounted = $state(false);

type Props = {
    data: SvelteMap<Date, number>;
};
let { data }: Props = $props();

let filteredData = $derived.by(() => {
    const transformed = data
        .entries()
        .map(([k, v]) => {
            return { date: k, value: v };
        })
        .toArray();

    transformed.sort((a, b) => (a.date > b.date ? 1 : -1));
    return transformed.slice(-MAX_RECORDS, transformed.length);
});

$effect(() => {
    if (!mounted) return;
    // TODO: yikes + do not shift if there is less than MAX_RECORDS
    const newData = filteredData.pop();
    chart.data?.labels?.shift();
    chart.data?.labels?.push(newData?.date.toUTCString());
    chart.data?.datasets[0].data.shift();
    chart.data?.datasets[0].data.push(newData?.value ?? 0);
    chart.update();
});

onMount(() => {
    chart = new Chart(ctx, {
        type: 'line',
        options: {
            plugins: {
                legend: {
                    display: false,
                },
                tooltip: {},
            },
            scales: {
                x: {
                    ticks: {
                        display: false,
                    },
                },
            },
        },
        data: {
            labels: filteredData.map((e) => e.date.toUTCString()),
            datasets: [
                {
                    label: '',
                    data: filteredData.map((e) => e.value),
                },
            ],
        },
    });

    mounted = true;
});
</script>

<canvas bind:this={ctx}></canvas>
