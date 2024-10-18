<script lang="ts">
    import { Label } from '$lib/components/ui/label';
    import * as Select from '$lib/components/ui/select';
    import { Button } from '$lib/components/ui/button';
    import Input from '../ui/input/input.svelte';
    import type { Sensor } from '@/types/sensor';
    import type {
        RuleAndType,
        RuleDetails,
        RuleInternal,
        RuleNotType,
        RuleOrType,
    } from '@/types/rule';
    import RuleInternalBuilder from './RuleInternalBuilder.svelte';

    export let open: boolean;
    export let sensors: Sensor[];
    export let parent: RuleDetails | RuleNotType | RuleAndType | RuleOrType;

    function isRuleDetails(
        parentInput: RuleInternal | RuleDetails
    ): parentInput is RuleDetails {
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

    function constructRule(): RuleInternal | undefined {
        if (selectedType.value === 'gt' || selectedType.value === 'lt') {
            return {
                type: selectedType.value,
                sensor_id: selectedSensor.value,
                value: value,
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
        if (isRuleDetails(parent)) {
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
</script>

{#if open}
    <div class="flex items-center gap-3 min-w-[32rem]">
        <Label>Type:</Label>

        <Select.Root bind:selected={selectedType} required name="type">
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
            <Select.Root bind:selected={selectedSensor} required name="type">
                <Select.Trigger>
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

            <Input type="number" bind:value />
        {/if}

        <div class="flex">
            <Button on:click={() => (open = false)} variant="outline">Cancel</Button>

            <Button on:click={addRule} variant="outline">Add</Button>
        </div>
    </div>
{/if}
