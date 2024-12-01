<script lang="ts">
    import { Chart } from 'chart.js/auto';
    import { onMount } from 'svelte';
    import type { SvelteMap } from 'svelte/reactivity';

    const DEFAULT_RECORD_COUNT = 32;

    let ctx: HTMLCanvasElement;
    let chart: Chart;
    let mounted = $state(false);

    type Props = {
        data: SvelteMap<Date, number>;
        fixedView: {
            from: Date;
            to: Date;
        } | null;
        accuracy?: number; // 0-1
        defaultRecordCount?: number;
    };
    let {
        data,
        fixedView = $bindable(null),
        accuracy = $bindable(0),
        defaultRecordCount = DEFAULT_RECORD_COUNT,
    }: Props = $props();

    let filteredData = $derived.by(() => {
        const transformed = data
            .entries()
            .map(([k, v]) => {
                return { date: k, value: v };
            })
            .toArray();

        transformed.sort((a, b) => (a.date > b.date ? 1 : -1));

        if (fixedView === null)
            return transformed.slice(-defaultRecordCount, transformed.length);

        const filtered = transformed.filter(
            ({ date }) => fixedView.from <= date && date <= fixedView.to
        );

        // Accuracy - record count
        // 0.0 - defaultRecordCount
        // 1.0 - filtered.length
        const maxEntries = Math.floor(
            defaultRecordCount +
                accuracy *
                    accuracy *
                    accuracy *
                    (filtered.length - defaultRecordCount)
        );

        // Aggregate data if it exceeds maxEntries
        if (filtered.length <= maxEntries || accuracy >= 1.0) return filtered;

        const compartmentSize = Math.ceil(filtered.length / maxEntries);
        const aggregated = [];
        for (let i = 0; i < filtered.length; i += compartmentSize) {
            const compartment = filtered.slice(i, i + compartmentSize);
            const avgDate = new Date(
                compartment.reduce(
                    (sum, item) => sum + item.date.getTime(),
                    0
                ) / compartment.length
            );
            const avgValue =
                compartment.reduce((sum, item) => sum + item.value, 0) /
                compartment.length;
            aggregated.push({ date: avgDate, value: avgValue });
        }
        return aggregated;
    });

    $effect(() => {
        if (!mounted) return;

        if (fixedView) {
            if (chart.options.scales?.x?.display !== undefined)
                chart.options.scales.x.display = true;
            const diffInHours =
                (filteredData[filteredData.length - 1]?.date.getTime() -
                    filteredData[0]?.date.getTime()) /
                (1000 * 60 * 60);
            const diffInDays = diffInHours / 24;

            const formatDate = (date: Date) => {
                // If the range is within a few minutes, show minutes and seconds
                if (diffInHours < 1) {
                    return new Date(date).toLocaleTimeString([], {
                        hour: '2-digit',
                        minute: '2-digit',
                        second: '2-digit',
                    });
                }

                // If the range is within a few hours, show only hours and minutes
                if (diffInHours < 24) {
                    return new Date(date).toLocaleTimeString([], {
                        hour: '2-digit',
                        minute: '2-digit',
                    });
                }

                // If the range is within a few days, show only the date (month/day)
                if (diffInDays < 7) {
                    return new Date(date).toLocaleDateString([], {
                        month: '2-digit',
                        day: '2-digit',
                    });
                }

                // Otherwise, return the full UTC date string
                return new Date(date).toUTCString();
            };

            chart.data.labels = [
                ...filteredData.values().map((e) => formatDate(e.date)),
            ];
            chart.data.datasets[0].data = [
                ...filteredData.values().map((e) => e.value),
            ];
        } else {
            if (chart.options.scales?.x?.display)
                chart.options.scales.x.display = false;
            if (
                chart.data?.labels?.[defaultRecordCount - 1] !==
                filteredData[filteredData.length - 2]?.date.toUTCString()
            ) {
                // still has fixed view data shown, have to clear it first
                if (chart.data?.labels)
                    chart.data.labels = [
                        ...filteredData
                            .values()
                            .map((e) => e.date.toUTCString()),
                    ];
                chart.data.datasets[0].data = [
                    ...filteredData.values().map((e) => e.value),
                ];
            }
            // TODO: yikes + do not shift if there is less than MAX_RECORDS
            const newData = filteredData.pop();
            chart.data?.labels?.shift();
            chart.data?.labels?.push(newData?.date.toUTCString());
            chart.data?.datasets[0].data.shift();
            chart.data?.datasets[0].data.push(newData?.value ?? 0);
        }

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
                            display: true,
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
