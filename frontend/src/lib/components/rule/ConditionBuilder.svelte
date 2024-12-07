<script lang="ts">
    import { Label } from '$lib/components/ui/label';
    import * as Select from '$lib/components/ui/select';
    import { Button } from '$lib/components/ui/button';
    import Input from '../ui/input/input.svelte';
    import type { Sensor } from '@/types/sensor';
    import type {
        NewRule,
        RuleAndType,
        RuleDetails,
        RuleInternal,
        RuleNotType,
        RuleOrType,
    } from '@/types/rule';

    type Props = {
        open: boolean;
        sensors: Sensor[];
        parent: RuleDetails | NewRule | RuleNotType | RuleAndType | RuleOrType;
    };

    let { open = $bindable(), sensors, parent = $bindable() }: Props = $props();

    function isRootRule(
        parentInput: RuleInternal | RuleDetails | NewRule
    ): parentInput is RuleDetails | NewRule {
        return Object.hasOwn(parentInput, 'description');
    }

    let selectedType = $state({
        value: 'and',
        label: 'And',
    });

    const typeConfig = [
        { value: 'gt', label: 'Greater than' },
        { value: 'lt', label: 'Lower than' },
        { value: 'and', label: 'And' },
        { value: 'or', label: 'Or' },
        { value: 'perc', label: 'Perc' },
        { value: 'time', label: 'Time' },
    ];

    let selectedVariant: { value: string; label: string } = $state({
        value: '',
        label: '---',
    });
    let selectedSensor: { value: string; label: string } = $state({
        value: '',
        label: '---',
    });
    let time: string = $state('-:-');
    let value: number = $state(0);
    let percentile: number = $state(0);
    let duration = $state({
        hours: '0',
        minutes: '0',
        seconds: '0',
    });
    let errors = $state({
        value: false,
        sensor: false,
        percentile: false,
        duration: {
            hours: false,
            minutes: false,
            seconds: false,
        },
        variant: false,
        time: false,
    });

    function constructRule(): RuleInternal | undefined {
        if (selectedType.value === 'gt' || selectedType.value === 'lt') {
            if (errors.value || errors.sensor) return;
            return {
                type: selectedType.value,
                sensor_id: selectedSensor.value,
                value: +value,
            };
        } else if (
            selectedType.value === 'and' ||
            selectedType.value === 'or'
        ) {
            return {
                type: selectedType.value,
                children: [],
            };
        } else if (selectedType.value === 'perc') {
            if (
                errors.duration.hours ||
                errors.duration.minutes ||
                errors.duration.seconds ||
                errors.percentile ||
                errors.percentile
            )
                return;
            return {
                type: selectedType.value,
                sensor_id: selectedSensor.value,
                duration: `${duration.hours}h${duration.minutes}m${duration.seconds}s`,
                perc: percentile,
            };
        } else if (selectedType.value === 'time') {
            if (errors.variant || errors.time) return;
            const match = time.split(':');
            const hour = parseInt(match[0], 10); // Extract hour (index 1)
            const minute = parseInt(match[1], 10); // Extract minute (index 2)

            return {
                type: selectedType.value,
                variant: selectedVariant.value as 'before' | 'after',
                hour,
                minute,
            };
        } else {
            return;
        }
    }

    function addRule() {
        const rule = constructRule();
        if (!rule) return;
        open = false;
        if (isRootRule(parent)) {
            parent.internal = rule;
            return;
        } else if (parent.type === 'or' || parent.type === 'and') {
            parent.children = [...parent.children, rule];
            return;
        } else if (parent.type === 'not') {
            parent.wrapped = rule;
            return;
        }
    }

    $effect(() => {
        errors.value = typeof value === 'undefined';
        errors.sensor = !selectedSensor;

        if (selectedType.value === 'perc') {
            errors.percentile = Number(percentile) <= 0;
            errors.sensor = !sensors.find((s) => s.id === selectedSensor.value);
            errors.duration.hours = Number(duration.hours) < 0;
            errors.duration.minutes = Number(duration.minutes) < 0;
            errors.duration.seconds = Number(duration.seconds) < 0;

            if (
                Number(duration.hours) +
                    Number(duration.minutes) +
                    Number(duration.seconds) <=
                0
            ) {
                errors.duration.hours = true;
                errors.duration.minutes = true;
                errors.duration.seconds = true;
            }
        } else if (selectedType.value === 'gt' || selectedType.value === 'lt') {
            errors.sensor = !sensors.find((s) => s.id === selectedSensor.value);
        } else if (selectedType.value === 'time') {
            errors.variant =
                selectedVariant.value !== 'before' &&
                selectedVariant.value !== 'after';
            const timeRegex = /^([01]\d|2[0-3]):[0-5]\d$/;
            errors.time = !time || !timeRegex.test(time);
        }
    });
</script>

{#if open}
    <div class="flex items-center gap-3">
        <Label>Type:</Label>

        <Select.Root bind:selected={selectedType}>
            <Select.Trigger>
                <Select.Value />
            </Select.Trigger>
            <Select.Content>
                {#each typeConfig as config}
                    <Select.Item value={config.value}
                        >{config.label}</Select.Item
                    >
                {/each}
            </Select.Content>
        </Select.Root>

        {#if selectedType.value === 'gt' || selectedType.value === 'lt'}
            <Label>Sensor:</Label>
            <Select.Root bind:selected={selectedSensor}>
                <Select.Trigger
                    class={errors['sensor'] ? 'border-2 border-red-600' : ''}
                >
                    <Select.Value />
                </Select.Trigger>
                <Select.Content>
                    {#each sensors as sensor}
                        <Select.Item value={sensor.id}
                            >{sensor.name}</Select.Item
                        >
                    {/each}
                </Select.Content>
            </Select.Root>

            <Label>value:</Label>

            <Input
                type="number"
                class="min-w-[4rem] {errors['value']
                    ? 'border-2 border-red-600'
                    : ''}"
                bind:value
            />
        {:else if selectedType.value === 'perc'}
            <Label>Sensor:</Label>
            <Select.Root bind:selected={selectedSensor}>
                <Select.Trigger
                    class={errors['sensor'] ? 'border-2 border-red-600' : ''}
                >
                    <Select.Value />
                </Select.Trigger>
                <Select.Content>
                    {#each sensors as sensor}
                        <Select.Item value={sensor.id}
                            >{sensor.name}</Select.Item
                        >
                    {/each}
                </Select.Content>
            </Select.Root>

            <Label>Percentile:</Label>
            <Input
                type="number"
                class="min-w-[4rem] {errors.percentile
                    ? 'border-2 border-red-600'
                    : ''}"
                bind:value={percentile}
            />

            <Label>For:</Label>
            <div class="flex items-center gap-2">
                <input
                    type="number"
                    class="time-part-input {errors.duration.hours
                        ? 'border-2 border-red-600'
                        : ''}"
                    bind:value={duration.hours}
                    min="0"
                    max="23"
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
                    placeholder="MM"
                />
                <span>:</span>
                <input
                    type="number"
                    class="time-part-input {errors.duration.seconds
                        ? 'border-2 border-red-600'
                        : ''}"
                    bind:value={duration.seconds}
                    min="0"
                    max="59"
                    placeholder="SS"
                />
                <Label class="">HH:MM:SS</Label>
            </div>
        {:else if selectedType.value === 'time'}
            <Label>Variant:</Label>
            <Select.Root bind:selected={selectedVariant}>
                <Select.Trigger
                    class={errors.variant ? 'border-2 border-red-600' : ''}
                >
                    <Select.Value />
                </Select.Trigger>
                <Select.Content>
                    <Select.Item value={'before'}>{'Before'}</Select.Item>
                    <Select.Item value={'after'}>{'After'}</Select.Item>
                </Select.Content>
            </Select.Root>

            <Input
                type="time"
                class="min-w-[4rem] {errors.time
                    ? 'border-2 border-red-600'
                    : ''}"
                bind:value={time}
            />
        {/if}

        <div class="flex">
            <Button on:click={() => (open = false)} variant="outline"
                >Cancel</Button
            >

            <Button on:click={addRule} variant="outline">Add</Button>
        </div>
    </div>
{/if}

<style>
    .time-part-input {
        text-align: center;
        width: 2.5rem;
    }
</style>
