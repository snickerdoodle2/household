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

    export let open: boolean;
    export let sensors: Sensor[];
    export let parent:
        | RuleDetails
        | NewRule
        | RuleNotType
        | RuleAndType
        | RuleOrType;

    function isRootRule(
        parentInput: RuleInternal | RuleDetails | NewRule
    ): parentInput is RuleDetails | NewRule {
        return Object.hasOwn(parentInput, 'description');
    }

    let selectedType = {
        value: 'and',
        label: 'And',
    };

    const typeConfig = [
        { value: 'gt', label: 'Greater than' },
        { value: 'lt', label: 'Lower than' },
        { value: 'and', label: 'And' },
        { value: 'or', label: 'Or' },
    ];

    let selectedSensor: { value: string; label: string };
    let value: number;
    let errors = { value: false, sensor: false };

    function constructRule(): RuleInternal | undefined {
        if (selectedType.value === 'gt' || selectedType.value === 'lt') {
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
        } else {
            return undefined;
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

    const validate = () => {
        errors.value = typeof value === 'undefined';
        errors.sensor = !selectedSensor;
    };

    let timeout: number;
    // eslint-disable-next-line @typescript-eslint/no-unsafe-function-type
    const debounce = (callback: Function, ...args: unknown[]) => {
        clearTimeout(timeout);
        timeout = window.setTimeout(() => callback(args), 300);
    };
    $: debounce(validate, selectedSensor, value);
</script>

{#if open}
    <div class="flex min-w-[35rem] items-center gap-3">
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
        {/if}

        <div class="flex">
            <Button on:click={() => (open = false)} variant="outline"
                >Cancel</Button
            >

            <Button on:click={addRule} variant="outline">Add</Button>
        </div>
    </div>
{/if}
