<script lang="ts">
    import { Input } from '@/components/ui/input';
    import type { RuleTimeType } from '@/types/rule';
    import type { Sensor } from '@/types/sensor';
    import * as Select from '$lib/components/ui/select';
    import { onMount } from 'svelte';
    import { Label } from '$lib/components/ui/label';
    import Button from '../ui/button/button.svelte';
    import { Disc, Pencil1 } from 'radix-icons-svelte';

    type Props = {
        internal: RuleTimeType;
        sensors: Sensor[];
        editingDisabled?: boolean;
        children?: import('svelte').Snippet;
    };

    let {
        internal = $bindable(),
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        sensors,
        editingDisabled = $bindable(false),
        children,
    }: Props = $props();

    let editing = $state(false);

    let selectedVariant: { value: string; label: string } = $state({
        value: '',
        label: '---',
    });
    let time: string = $state('-:-');

    let dropDownsOpen = $state({
        sensor: false,
        type: false,
    });
    let errors = $state({
        variant: false,
        time: false,
    });

    let wrappingDiv: HTMLDivElement | undefined = $state();

    function toggleEditing() {
        editing = !editing;
    }

    function capitalizeFirstLetter(str: string) {
        if (!str) return str; // Handle empty strings
        return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase();
    }

    function save() {
        errors.variant =
            selectedVariant.value !== 'before' &&
            selectedVariant.value !== 'after';
        const timeRegex = /^([01]\d|2[0-3]):[0-5]\d$/;
        errors.time = !time || !timeRegex.test(time);

        if (errors.variant || errors.time) return false;

        const match = time.split(':');
        const hour = parseInt(match[0], 10); // Extract hour (index 1)
        const minute = parseInt(match[1], 10); // Extract minute (index 2)

        internal.hour = hour;
        internal.minute = minute;
        internal.variant = selectedVariant.value as 'before' | 'after';
        return true;
    }

    function syncInternalValues() {
        time = `${internal.hour}:${internal.minute}`;
        selectedVariant = {
            label: capitalizeFirstLetter(internal.variant),
            value: internal.variant,
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
    <Label>Time</Label>

    <Select.Root
        bind:selected={selectedVariant}
        bind:open={dropDownsOpen.sensor}
        required
        disabled={!editing}
    >
        <Select.Trigger class={errors.variant ? 'border-2 border-red-600' : ''}>
            <Select.Value />
        </Select.Trigger>
        <Select.Content>
            <Select.Item value={'before'}>{'Before'}</Select.Item>
            <Select.Item value={'after'}>{'After'}</Select.Item>
        </Select.Content>
    </Select.Root>

    <Input
        type="time"
        class="min-w-[4rem] {errors.time ? 'border-2 border-red-600' : ''}"
        disabled={!editing}
        bind:value={time}
    />

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
