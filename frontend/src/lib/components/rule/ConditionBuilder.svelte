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
    ];

    let selectedSensor: { value: string; label: string } = $state({
        value: '',
        label: '',
    });
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
            <div class="flex flex-row items-center rounded-md gap-1">
                <div
                    class="flex items-center {errors.duration.hours
                        ? 'border-2 border-red-600'
                        : ''}"
                >
                    <Input
                        type="number"
                        class="min-w-14"
                        bind:value={duration.hours}
                    />
                    <Label>h</Label>
                </div>

                <div
                    class="flex items-center {errors.duration.minutes
                        ? 'border-2 border-red-600'
                        : ''}"
                >
                    <Input
                        type="number"
                        class="min-w-14"
                        bind:value={duration.minutes}
                    />
                    <Label>m</Label>
                </div>

                <div
                    class="flex items-center {errors.duration.seconds
                        ? 'border-2 border-red-600'
                        : ''}"
                >
                    <Input
                        type="number"
                        class="min-w-14"
                        bind:value={duration.seconds}
                    />
                    <Label>s</Label>
                </div>
            </div>
        {/if}

        <div class="flex">
            <Button on:click={() => (open = false)} variant="outline"
                >Cancel</Button
            >

            <Button on:click={addRule} variant="outline">Add</Button>
        </div>
    </div>
{/if}
