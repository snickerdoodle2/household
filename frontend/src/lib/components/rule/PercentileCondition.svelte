<script lang="ts">
    import { Input } from '@/components/ui/input';
    import type { RulePercType } from '@/types/rule';
    import type { Sensor } from '@/types/sensor';
    import * as Select from '$lib/components/ui/select';
    import { onMount } from 'svelte';
    import { Label } from '$lib/components/ui/label';
    import Button from '../ui/button/button.svelte';
    import { Disc, Pencil1 } from 'radix-icons-svelte';
    import { parseDuration } from '@/helpers/time';

    type Props = {
        internal: RulePercType;
        sensors: Sensor[];
        editingDisabled?: boolean;
        children?: import('svelte').Snippet;
    };

    let {
        internal = $bindable(),
        sensors,
        editingDisabled = $bindable(false),
        children,
    }: Props = $props();

    let editing = $state(false);

    let selectedSensor: { value: string; label: string } = $state({
        value: '',
        label: '',
    });
    let dropDownsOpen = $state({
        sensor: false,
        type: false,
    });
    let percentile: string = $state('0');
    let duration = $state({
        hours: '0',
        minutes: '0',
        seconds: '0',
    });

    let errors = $state({
        percentile: false,
        duration: {
            hours: false,
            minutes: false,
            seconds: false,
        },
        sensor_id: false,
    });

    let wrappingDiv: HTMLDivElement | undefined = $state();

    function toggleEditing() {
        editing = !editing;
    }

    function save() {
        errors.percentile = Number(percentile) <= 0;
        errors.sensor_id = !sensors.find((s) => s.id === selectedSensor.value);
        errors.duration.hours = Number(duration.hours) < 0;
        errors.duration.minutes = Number(duration.minutes) < 0;
        errors.duration.seconds = Number(duration.seconds) < 0;

        if (
            Number(duration.minutes) +
                Number(duration.minutes) +
                Number(duration.seconds) <=
            0
        ) {
            errors.duration.hours = true;
            errors.duration.minutes = true;
            errors.duration.seconds = true;
        }

        if (
            errors.percentile ||
            errors.sensor_id ||
            errors.duration.hours ||
            errors.duration.minutes ||
            errors.duration.seconds
        )
            return false;

        internal.perc = Number(percentile);
        internal.duration = `${duration.hours}h${duration.minutes}m${duration.seconds}s`;
        internal.sensor_id = selectedSensor.value;

        return true;
    }

    function syncInternalValues() {
        const initialSensor = sensors.find(
            (sensor) => sensor.id === internal.sensor_id
        );
        selectedSensor = {
            value: initialSensor?.id ?? 'not found',
            label: initialSensor?.name ?? 'not found',
        };
        percentile = internal.perc.toString();
        const parsedDuration = parseDuration(internal.duration);
        duration = {
            hours: parsedDuration.hours.toString(),
            minutes: parsedDuration.minutes.toString(),
            seconds: parsedDuration.seconds.toString(),
        };
    }

    function handleClick(event: MouseEvent) {
        if (
            wrappingDiv &&
            !wrappingDiv.contains(event.target as Node) &&
            !event.defaultPrevented &&
            !Object.values(dropDownsOpen).some((open) => open)
        ) {
            editing = false;
            syncInternalValues();
        }
    }

    onMount(() => {
        // initialize values
        syncInternalValues();

        // Capture phase to ensure clicks inside are checked before propagation finishes
        document.addEventListener('pointerdown', handleClick, true);
    });
</script>

<div
    class="flex w-full flex-row items-center gap-2 whitespace-nowrap"
    bind:this={wrappingDiv}
>
    <Label>Value of</Label>

    <Select.Root
        bind:selected={selectedSensor}
        bind:open={dropDownsOpen.sensor}
        required
        name="type"
        disabled={!editing}
    >
        <Select.Trigger
            class={errors.sensor_id ? 'border-2 border-red-600' : ''}
        >
            <Select.Value />
        </Select.Trigger>
        <Select.Content>
            {#each sensors as sensor}
                <Select.Item value={sensor.id}>{sensor.name}</Select.Item>
            {/each}
        </Select.Content>
    </Select.Root>

    <Label>Percentile:</Label>
    <Input
        type="number"
        class="min-w-14 {errors.percentile ? 'border-2 border-red-600' : ''}"
        bind:value={percentile}
        disabled={!editing}
    />

    <Label>For:</Label>
    <div
        class="flex items-center gap-2"
    >
        <input
            type="number"
            class="time-part-input {errors.duration.hours
                ? 'border-2 border-red-600'
                : ''}"
            bind:value={duration.hours}
            min="0"
            max="23"
            disabled={!editing}
            placeholder="HH"
        />
        <span>:</span>
        <input
            type="number"
            class="time-part-input {errors.duration.minutes
                ? 'border-2 border-red-600'
                : ''}"
            bind:value={duration.minutes}
            min="0"
            max="59"
            disabled={!editing}
            placeholder="MM"
        />
        <span>:</span>
        <input
            type="number"
            class="time-part-input"
            bind:value={duration.seconds}
            min="0"
            max="59"
            disabled={!editing}
            placeholder="SS"
        />
    </div>
    <Label class="">HH:MM:SS</Label>

    {#if !editingDisabled}
        {#if editing}
            <Button
                variant="outline"
                size="icon"
                on:click={() => {
                    if (save()) toggleEditing();
                }}
            >
                <Disc />
            </Button>
        {:else}
            <Button variant="outline" size="icon" on:click={toggleEditing}>
                <Pencil1 />
            </Button>
        {/if}
    {/if}

    {@render children?.()}
</div>

<style>
    .time-part-input {
        text-align: center;
        width: 2.5rem;
    }
</style>
