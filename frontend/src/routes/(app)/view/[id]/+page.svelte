<script lang="ts">
    import type { PageData } from './$types';
    import { onMount } from 'svelte';
    import * as Dialog from '$lib/components/ui/dialog';
    import * as Card from '$lib/components/ui/card';
    import { goto } from '$app/navigation';
    import { SensorWebsocket } from '@/helpers/socket.svelte';
    import type { SvelteMap } from 'svelte/reactivity';
    import Chart from '@/components/sensor/Chart.svelte';
    import Button from '@/components/ui/button/button.svelte';
    import type { SensorDetails } from '@/types/sensor';
    import { Label } from '$lib/components/ui/label';

    const DEFAULT_RECORD_COUNT = 32;

    const StatisticsTypes = {
        Mean: 'mean',
        Median: 'median',
        Min: 'min',
        Max: 'max',
        StandardDeviation: 'standardDeviation',
        Trend: 'trend',
    } as const;

    const STATISTIC_NAMES = {
        [StatisticsTypes.Mean]: 'Mean',
        [StatisticsTypes.Median]: 'Median',
        [StatisticsTypes.Min]: 'Min',
        [StatisticsTypes.Max]: 'Max',
        [StatisticsTypes.StandardDeviation]: 'Standard Deviation',
        [StatisticsTypes.Trend]: 'Trend',
    };

    type Props = {
        data: PageData;
    };

    let { data }: Props = $props();

    let sensorId: string = $state(data.sensorId);
    let sensor: SensorDetails | undefined = $state(undefined);
    let updateValues = $state(false);
    let fixedView: {
        from: Date;
        to: Date;
    } | null = $state(null);
    let accuracy = $state(0.0)

    let statistics = $state({
        [StatisticsTypes.Mean]: 0,
        [StatisticsTypes.Median]: 0,
        [StatisticsTypes.Min]: 0,
        [StatisticsTypes.Max]: 0,
        [StatisticsTypes.StandardDeviation]: 0,
        [StatisticsTypes.Trend]: 0,
    });

    onMount(async () => {
        sensor = await data.sensor;
    });

    const close = () => {
        goto(`/`);
    };

    $effect(() => {
        const entries = [...(ws.data.get(sensorId)?.entries() ?? [])];
        if (!entries.length) return;

        let data: number[] = [];
        if (fixedView !== null) {
            const { from, to } = fixedView;
            data = entries
                .filter(([d]) => from <= d && d <= to)
                .map(([, v]) => v);
        } else {
            data = entries
                .map(([, v]) => v)
                .slice(-DEFAULT_RECORD_COUNT, entries.length);
        }

        // If no data, reset statistics
        if (!data.length) {
            statistics = {
                mean: 0,
                median: 0,
                min: 0,
                max: 0,
                standardDeviation: 0,
                trend: 0,
            };
            return;
        }

        // Calculate Mean
        const mean = data.reduce((sum, value) => sum + value, 0) / data.length;

        // Calculate Median
        const sorted = [...data].sort((a, b) => a - b);
        const mid = Math.floor(sorted.length / 2);
        const median =
            sorted.length % 2 === 0
                ? (sorted[mid - 1] + sorted[mid]) / 2
                : sorted[mid];

        // Calculate Min and Max
        const min = Math.min(...data);
        const max = Math.max(...data);

        // Calculate Standard Deviation
        const variance =
            data.reduce((sum, value) => sum + Math.pow(value - mean, 2), 0) /
            data.length;
        const standardDeviation = Math.sqrt(variance);

        // Calculate Trend (Simple Linear Regression)
        const n = data.length;
        const x = entries.slice(-n).map(([d]) => new Date(d).getTime()); // x-axis: timestamps
        const y = data; // y-axis: sensor values
        const xSum = x.reduce((sum, value) => sum + value, 0);
        const ySum = y.reduce((sum, value) => sum + value, 0);
        const xySum = x.reduce((sum, value, i) => sum + value * y[i], 0);
        const xSquaredSum = x.reduce((sum, value) => sum + value * value, 0);
        const trend = (n * xSquaredSum - Math.pow(xSum, 2) == 0) ? 0 : (n * xySum - xSum * ySum) / (n * xSquaredSum - Math.pow(xSum, 2)); // Represents the rate of change over time

        // Update the statistics state
        statistics = {
            mean,
            median,
            min,
            max,
            standardDeviation,
            trend,
        };
    });

    // websocket stuff
    const ws = new SensorWebsocket();
    $effect(() => {
        if (!ws.ready) return;
        ws.subscribe(sensorId);

        return () => {
            ws.unsubscribe(sensorId);
        };
    });

    let chartData: SvelteMap<Date, number> | undefined = $derived(
        ws.data.get(sensorId)
    );
</script>

<Dialog.Root
    open={true}
    onOpenChange={(opened) => {
        if (!opened) {
            close();
        }
    }}
>
    <Dialog.Portal>
        <Dialog.Overlay />
        <Dialog.Content
            class="flex flex-row max-w-none items-center justify-center px-8 py-4 md:w-fit"
        >
            {#if !sensor}
                <p>Loading...</p>
            {:else}
                <Card.Root class="min-w-[800px] border-none shadow-none">
                    <Card.Header class="text-center text-3xl">
                        <Card.Title>
                            {sensor.name}
                            <span class="text-xl text-gray-500 align-center">
                                ({sensor.type})
                            </span>
                        </Card.Title>
                    </Card.Header>
                    <Card.Content>
                        <div class="flex w-full">
                            <div class="p-2 w-full">
                                {#if chartData}
                                    <Chart
                                        data={chartData}
                                        bind:fixedView
                                        bind:accuracy
                                        defaultRecordCount={DEFAULT_RECORD_COUNT}
                                    />
                                {:else}
                                    <p>Error opening socket</p>
                                {/if}
                            </div>
                            <div>
                                <table class="ml-4 text-sm">
                                    <tbody>
                                        {#each Object.values(StatisticsTypes) as statisticType, idx}
                                            <tr
                                                class="border-b {idx == 0
                                                    ? 'border-t'
                                                    : ''}"
                                            >
                                                <td class="px-4 py-2"
                                                    >{STATISTIC_NAMES[
                                                        statisticType
                                                    ]}</td
                                                >
                                                <td
                                                    class="px-4 py-2 text-center min-w-16"
                                                    >{statistics[
                                                        statisticType
                                                    ].toFixed(2)}</td
                                                >
                                            </tr>
                                        {/each}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                        <div>
                            <div
                                class="flex items-center align-center gap-2 min-w-32"
                            >
                                <Label class="mr-2">Accuracy:</Label>

                                <input
                                    type="range"
                                    class="mt-1"
                                    min="0.0"
                                    max="1.0"
                                    step="0.05"
                                    bind:value={accuracy}
                                    disabled={!fixedView}
                                />

                                <Label>
                                    <span class="text-md text-gray-500">
                                        {accuracy.toFixed(2)}
                                    </span>
                                </Label>
                            </div>
                        </div>
                    </Card.Content>
                    <Card.Footer class="flex justify-between">
                        <div>
                            {#if fixedView}
                                <Button
                                    on:click={() => {
                                        updateValues = false;
                                        fixedView = null;
                                    }}
                                    size="bold">Resume</Button
                                >
                            {/if}

                            <Button
                                on:click={() => {
                                    updateValues = false;

                                    ws.requestSince(sensorId, '24h');
                                    fixedView = {
                                        from: new Date(
                                            Date.now() - 60 * 60 * 1000
                                        ), // One hour ago
                                        to: new Date(), // Current time
                                    };
                                }}
                                size="bold">New</Button
                            >
                        </div>
                        <div>
                            <Button
                                on:click={() => {
                                    goto(`/details/${sensorId}`);
                                }}
                                size="bold">Go to details</Button
                            >
                        </div>
                    </Card.Footer>
                </Card.Root>
            {/if}
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>
