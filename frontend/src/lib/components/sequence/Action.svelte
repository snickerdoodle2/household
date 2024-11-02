<script lang="ts">
import type { SequenceAction } from '@/types/sequence';
import * as Select from '../ui/select';
import { onMount } from 'svelte';
import Input from '../ui/input/input.svelte';
import Label from '../ui/label/label.svelte';

export let action: SequenceAction;
export let sensors: { label: string; value: string }[];
export let editing = false;
export let errorFields: string[] | undefined = [];

let selectedSensor = { value: 'unknown', label: 'unknown' };
let value = '0';
let time = {
    hours: 0,
    minutes: 0,
    seconds: 0,
};

$: if (action) syncActionValues();
$: action.value = Number(value);
$: action.msDelay = convertTimeToMs(time);
$: action.target = selectedSensor.value;

function convertMsToTime(ms: number): {
    hours: number;
    minutes: number;
    seconds: number;
} {
    const seconds = Math.floor(Math.max(0, ms / 1000) % 60);
    const minutes = Math.floor(Math.max(0, ms / (1000 * 60)) % 60);
    const hours = Math.floor(Math.max(0, ms / (1000 * 60 * 60)));

    return {
        hours,
        minutes,
        seconds,
    };
}

function convertTimeToMs(time: {
    hours: number;
    minutes: number;
    seconds: number;
}): number {
    const { hours, minutes, seconds } = time;

    const ms = hours * 60 * 60 * 1000 + minutes * 60 * 1000 + seconds * 1000;

    return ms;
}

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
