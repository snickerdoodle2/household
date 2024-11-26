<script lang="ts">
    import { Input } from '@/components/ui/input';
    import type { RuleGtType, RuleLtType } from '@/types/rule';
    import type { Sensor } from '@/types/sensor';
    import * as Select from '$lib/components/ui/select';
    import { onMount } from 'svelte';
    import { Label } from '$lib/components/ui/label';
    import Button from '../ui/button/button.svelte';
    import { Disc, Pencil1 } from 'radix-icons-svelte';

    type Props = {
        internal: RuleGtType | RuleLtType;
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

    let type: { value: 'gt' | 'lt' | 'not found'; label: string } = $state({
        value: 'gt',
        label: '',
    });
    let selectedSensor: { value: string; label: string } = $state({
        value: '',
        label: '',
    });
    let value: number = $state(1);
    let dropDownsOpen = $state({
        sensor: false,
        type: false,
    });

    let wrappingDiv: HTMLDivElement | undefined = $state();

    function toggleEditing() {
        editing = !editing;
    }

    function save() {
        if (type.value === 'not found') return;
        internal.type = type.value;
        internal.value = +value;
        internal.sensor_id = selectedSensor.value;
    }

    function syncInternalValues() {
        type = {
            value: internal.type,
            label: internal.type === 'gt' ? 'Greater than' : 'Lower than',
        };
        const initialSensor = sensors.find(
            (sensor) => sensor.id === internal.sensor_id
        );
        selectedSensor = {
            value: initialSensor?.id ?? 'not found',
            label: initialSensor?.name ?? 'not found',
        };
        value = +internal.value;
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
        <Select.Trigger>
            <Select.Value />
        </Select.Trigger>
        <Select.Content>
            {#each sensors as sensor}
                <Select.Item value={sensor.id}>{sensor.name}</Select.Item>
            {/each}
        </Select.Content>
    </Select.Root>

    <Label>is</Label>

    <div class="min-w-[8rem]">
        <Select.Root
            bind:selected={type}
            bind:open={dropDownsOpen.type}
            required
            name="type"
            disabled={!editing}
        >
            <Select.Trigger>
                <Select.Value />
            </Select.Trigger>
            <Select.Content>
                <Select.Item value={'lt'}>{'Lower than'}</Select.Item>
                <Select.Item value={'gt'}>{'Greater than'}</Select.Item>
            </Select.Content>
        </Select.Root>
    </div>

    <Input type="number" bind:value disabled={!editing} />

    {#if !editingDisabled}
        {#if editing}
            <Button
                variant="outline"
                size="icon"
                on:click={() => {
                    save();
                    toggleEditing();
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
