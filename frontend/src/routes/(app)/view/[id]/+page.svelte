<script lang="ts">
    import { generateDurationString } from '@/helpers/time';

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
    import 'flatpickr/dist/flatpickr.css';
    import { DateInput, DatePicker, localeFromDateFnsLocale } from 'date-picker-svelte'
    import Input from '@/components/ui/input/input.svelte';

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
    let settingView = $state(false);
    let settingValue = $state(false);
    let valueToSet = $state(0);
    let fixedView: {
        from: Date,
        to: Date
    } | null = $state(null);
    let accuracy = $state(0.0)
    let startDate: Date= $state(new Date());
    let endDate: Date= $state(new Date());

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

    const setSensorValue = () => {
        console.log(`Setting sensor value not implemented.`, sensorId, valueToSet);
    }

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
                            <div class="p-2 pb-4 w-full">
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

                        {#if settingValue}
                            <div class="flex gap-2 text-md p-2">
                                <div class="pr-2">
                                    <Label>Value</Label>
                                    <Input
                                        type="number"
                                        size="sm"
                                        class="ml-2"
                                        bind:value={valueToSet}
                                    />
                                </div>
                                <Button
                                    size="sm"
                                    class="mt-6"
                                    on:click={setSensorValue}>Set</Button
                                >
                                <Button
                                    size="sm"
                                    class="mt-6"
                                    variant="destructive"
                                    on:click={() => (settingValue = false)}
                                    >Cancel</Button
                                >
                            </div>
                        {:else}
                            <Button
                                on:click={() => {
                                    settingValue = true;
                                }}
                                size="bold">Set sensor value</Button
                            >
                        {/if}

                        {#if settingView || fixedView}
                            <div class="flex items-left">
                                <div class="p-2 flex gap-4 items-center">
                                    <div class="text-sm">
                                        <Label class="text-md"
                                            >Start date:</Label
                                        >
                                        <DateInput
                                            class="ml-2"
                                            bind:value={startDate}
                                            format="yyyy/MM/dd HH:mm:ss"
                                            placeholder="2000/31/12 23:59:59"
                                            dynamicPositioning={true}
                                            timePrecision={'second'}
                                            disabled={!!fixedView}
                                        />
                                    </div>

                                    <div class="text-sm">
                                        <Label class="text-md">End date:</Label>
                                        <DateInput
                                            class="ml-2"
                                            bind:value={endDate}
                                            format="yyyy/MM/dd HH:mm:ss"
                                            placeholder="2000/31/12 23:59:59"
                                            dynamicPositioning={true}
                                            timePrecision={'second'}
                                            disabled={!!fixedView}
                                        />
                                    </div>
                                    {#if !fixedView}
                                        <Button
                                            class="mt-5"
                                            size="sm"
                                            on:click={() => {
                                                updateValues = false;
                                                ws.requestSince(
                                                    sensorId,
                                                    generateDurationString(
                                                        startDate,
                                                        new Date()
                                                    )
                                                );
                                                fixedView = {
                                                    from: startDate,
                                                    to: endDate,
                                                };
                                                settingView = false;
                                            }}>Show measurments</Button
                                        >
                                        <Button
                                            class="mt-5"
                                            variant="destructive"
                                            size="sm"
                                            on:click={() => {
                                                settingView = false;
                                            }}>Cancel</Button
                                        >
                                    {:else}
                                        <div>
                                            <Label class="text-md"
                                                >Accuracy:</Label
                                            >
                                            <div class="ml-2 flex items-center">
                                                <div>
                                                    <input
                                                        type="range"
                                                        min="0.0"
                                                        max="1.0"
                                                        step="0.05"
                                                        bind:value={accuracy}
                                                        disabled={!fixedView}
                                                    />
                                                </div>

                                                <Label>
                                                    <span
                                                        class="pl-2 mb-4 text-md text-gray-500"
                                                    >
                                                        {accuracy.toFixed(2)}
                                                    </span>
                                                </Label>
                                            </div>
                                        </div>
                                        <Button
                                            class="mt-4"
                                            on:click={() => {
                                                updateValues = false;
                                                fixedView = null;
                                            }}
                                            size="sm">Resume</Button
                                        >
                                    {/if}
                                </div>
                            </div>
                        {:else if !fixedView && !settingView}
                            <Button
                                class="p-2"
                                on:click={() => {
                                    settingView = true;
                                }}
                                size="bold">See history</Button
                            >
                        {/if}
                    </Card.Content>
                    <Card.Footer class="flex justify-end">
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
