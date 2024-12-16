<script lang="ts">
    import * as Select from '$lib/components/ui/select';
    import Label from '../ui/label/label.svelte';

    type Props = {
        format: string | null;
        forcedFormat?: string;
        disabled?: boolean;
        dropDownsOpen?: boolean;
    };

    let {
        format = $bindable(),
        forcedFormat = $bindable(),
        disabled = $bindable(),
        dropDownsOpen: dropDownOpen = $bindable(),
    }: Props = $props();

    format = null;
    const months = [
        'January',
        'February',
        'March',
        'April',
        'May',
        'June',
        'July',
        'August',
        'September',
        'October',
        'November',
        'December',
    ];

    const weekdays = [
        'Monday',
        'Tuesday',
        'Wednesday',
        'Thursday',
        'Friday',
        'Saturday',
        'Sunday',
    ];

    let selectedDay: { value: string; label: string }[] = $state([]);
    let selectedMonth: { value: string; label: string }[] = $state([]);
    let selectedWeekDay: { value: string; label: string }[] = $state([]);

    let dayDropdownOpen = $state(false);
    let monthDropdownOpen = $state(false);
    let weekDayDropdownOpen = $state(false);

    $effect(() => {
        if (
            selectedDay.length === 0 ||
            selectedMonth.length === 0 ||
            selectedWeekDay.length === 0
        ) {
            format = null;
            return;
        }
        let dayString = selectedDay.find((e) => e.value == 'all')
            ? '*'
            : selectedDay.map((v) => v.value).join(',');
        let monthString = selectedMonth.find((e) => e.value == 'all')
            ? '*'
            : selectedMonth.map((v) => v.value).join(',');
        let weekDayString = selectedWeekDay.find((e) => e.value == 'all')
            ? '*'
            : selectedWeekDay.map((v) => v.value).join(',');

        format = `${dayString} ${monthString} ${weekDayString}`;
    });

    $effect(() => {
        dropDownOpen =
            dayDropdownOpen || monthDropdownOpen || weekDayDropdownOpen;
    });

    const expandRange = (range: string) => {
        const parts = range.split('-');
        if (
            parts.length === 2 &&
            !isNaN(Number(parts[0])) &&
            !isNaN(Number(parts[1]))
        ) {
            const start = parseInt(parts[0], 10);
            const end = parseInt(parts[1], 10);
            return Array.from({ length: end - start + 1 }, (_, i) =>
                (start + i).toString()
            );
        }
        return [range]; // Return the value as is if it's not a range
    };

    $effect(() => {
        if (!forcedFormat) return;
        // Reverse the initialFormat on mount and detect the selected days, months, and weekdays
        const [dayPart, monthPart, weekDayPart] = forcedFormat.split(' ');

        // Helper function to map part to its full text label
        const getDayOptions = (value: string) => {
            if (value === '*') return [{ value: 'all', label: 'All Days' }];
            return value
                .split(',')
                .flatMap((v) =>
                    expandRange(v).map((d) => ({
                        value: d.toString(),
                        label: `${d}`,
                    }))
                );
        };

        const getMonthOptions = (value: string) => {
            if (value === '*') return [{ value: 'all', label: 'All Months' }];
            return value
                .split(',')
                .flatMap((v) =>
                    expandRange(v).map((m) => ({
                        value: m.toString(),
                        label: months[parseInt(m)],
                    }))
                );
        };

        const getWeekdayOptions = (value: string) => {
            if (value === '*') return [{ value: 'all', label: 'All Weekdays' }];
            return value
                .split(',')
                .flatMap((v) =>
                    expandRange(v).map((w) => ({
                        value: w.toString(),
                        label: weekdays[parseInt(w)],
                    }))
                );
        };

        selectedDay = getDayOptions(dayPart);
        selectedMonth = getMonthOptions(monthPart);
        selectedWeekDay = getWeekdayOptions(weekDayPart);
        forcedFormat = undefined;
    });
</script>

<!-- Day number Select -->
<Label>Day</Label>
<Select.Root
    bind:selected={selectedDay}
    multiple
    {disabled}
    bind:open={dayDropdownOpen}
>
    <Select.Trigger
        class="min-w-16 {selectedDay.length == 0
            ? 'border-2 border-red-600'
            : ''}"
    >
        <Select.Value />
    </Select.Trigger>
    <Select.Content class="overflow-y-auto h-64">
        {#each Array(31).keys() as idx}
            <Select.Item value={idx.toString()}>{idx.toString()}</Select.Item>
        {/each}
        <Select.Item value={'all'}>{'All Days'}</Select.Item>
    </Select.Content>
</Select.Root>

<!-- Month Select -->
<Label>Month</Label>
<Select.Root
    bind:selected={selectedMonth}
    multiple
    {disabled}
    bind:open={monthDropdownOpen}
>
    <Select.Trigger
        class="min-w-32 {selectedMonth.length == 0
            ? 'border-2 border-red-600'
            : ''}"
    >
        <Select.Value />
    </Select.Trigger>
    <Select.Content>
        <Select.Item value={'all'}>{'All Months'}</Select.Item>
        {#each months as month, index}
            <Select.Item value={index.toString()}>{month}</Select.Item>
        {/each}
    </Select.Content>
</Select.Root>

<!-- Weekday Select -->
<Label>Weekday</Label>
<Select.Root
    bind:selected={selectedWeekDay}
    multiple
    {disabled}
    bind:open={weekDayDropdownOpen}
>
    <Select.Trigger
        class="min-w-32 {selectedWeekDay.length == 0
            ? 'border-2 border-red-600'
            : ''}"
    >
        <Select.Value />
    </Select.Trigger>
    <Select.Content>
        <Select.Item value={'all'}>{'All Weekdays'}</Select.Item>
        {#each weekdays as weekday, index}
            <Select.Item value={index.toString()}>{weekday}</Select.Item>
        {/each}
    </Select.Content>
</Select.Root>
