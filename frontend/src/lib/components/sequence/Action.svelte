<script lang="ts">
    import type { SequenceAction } from '@/types/sequence';
    import * as Select from '../ui/select';
    import { onMount } from 'svelte';
    import Input from '../ui/input/input.svelte';
    import Label from '../ui/label/label.svelte';
    import { convertMsToTime, convertTimeToMs } from '@/helpers/time';

    type Props = {
        action: SequenceAction;
        sensors: { label: string; value: string }[];
        editing: boolean;
        errorFields: string[] | undefined;
    };

    let {
        action = $bindable(),
        sensors = $bindable(),
        editing = $bindable(false),
        errorFields = $bindable([]),
    }: Props = $props();

    let selectedSensor = $state({ value: 'unknown', label: 'unknown' });
    let value = $state('0');
    let time = $state({
        hours: 0,
        minutes: 0,
        seconds: 0,
    });

    $effect(() => {
        if (action) syncActionValues();
    });
    $effect(() => {
        action.value = Number(value);
    });
    $effect(() => {
        action.msDelay = convertTimeToMs(time);
    });
    $effect(() => {
        action.target = selectedSensor.value;
    });

    function syncActionValues() {
        const initialSensor = sensors.find(
            (sensor) => sensor.value === action.target
        );
        selectedSensor = {
            label: initialSensor?.label ?? 'choose sensor',
            value: initialSensor?.value ?? 'choose sensor',
        };
        time = convertMsToTime(action.msDelay);
        value = action.value.toString();
    }

    onMount(() => {
        syncActionValues();
    });
</script>

<div
    class="flex w-full min-w-[40rem] flex-row items-center gap-2 whitespace-nowrap"
>
    <Label>Sensor:</Label>
    <Select.Root
        bind:selected={selectedSensor}
        required
        name="target"
        disabled={!editing}
    >
        <Select.Trigger
            class="min-w-[150px]{errorFields && errorFields.includes('target')
                ? 'border-2 border-red-600'
                : ''}"
        >
            <Select.Value />
        </Select.Trigger>
        <Select.Content>
            {#each sensors as sensor}
                <Select.Item value={sensor.value}>{sensor.label}</Select.Item>
            {/each}
        </Select.Content>
    </Select.Root>

    <Label class="pl-3">Value to send:</Label>
    <Input
        type="number"
        class={errorFields && errorFields.includes('value')
            ? 'border-2 border-red-600'
            : ''}
        bind:value
        disabled={!editing}
    />

    <Label class="pl-3">Delay:</Label>
    <div
        class="flex flex-row items-center rounded-md gap-1{errorFields &&
        errorFields.includes('msDelay')
            ? 'border-2 border-red-600'
            : ''}"
    >
        <div class="flex items-center">
            <Input type="number" bind:value={time.hours} disabled={!editing} />
            <Label>h</Label>
        </div>

        <div class="flex items-center">
            <Input
                type="number"
                bind:value={time.minutes}
                disabled={!editing}
            />
            <Label>m</Label>
        </div>

        <div class="flex items-center">
            <Input
                type="number"
                bind:value={time.seconds}
                disabled={!editing}
            />
            <Label>s</Label>
        </div>
    </div>
</div>
