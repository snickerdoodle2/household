<script lang="ts">
    import type { SequenceAction } from "@/types/sequence";
    import * as Select from "../ui/select";
    import { onMount } from "svelte";
    import Input from "../ui/input/input.svelte";
    import Label from "../ui/label/label.svelte";
    
    export let action: SequenceAction;
    export let sensors: { label: string; value: string }[]
    export let editing = false;

    let selectedSensor: { value: string; label: string };
    let value: string;
    let minuteDelay: string;
    
    $: action && syncActionValues();
    
    function syncActionValues() {
        const initialSensor = sensors.find(
            (sensor) => sensor.value === action.target
        );
        selectedSensor = {
            label: initialSensor?.label ?? 'choose sensor',
            value: initialSensor?.value ?? 'choose sensor',
        };
        minuteDelay = (action.msDelay / 1000 / 60).toFixed(2).toString();
        value = action.value.toString()
    }

    onMount(() => {
        syncActionValues()
    })


    // $: {
    //     action.value = Number(value);
    //     action.msDelay = Number(minuteDelay) * 1000 * 60;
    //     action.target = selectedSensor.value;
    // }
</script>

<div
    class="flex w-full min-w-[40rem] flex-row items-center gap-2 whitespace-nowrap"
>
    <Label>Sensor:</Label>
    <Select.Root
        bind:selected={selectedSensor}
        required
        name="Sensor"
        disabled={!editing}
    >
        <Select.Trigger>
            <Select.Value />
        </Select.Trigger>
        <Select.Content>
            {#each sensors as sensor}
                <Select.Item value={sensor.value}>{sensor.label}</Select.Item>
            {/each}
        </Select.Content>
    </Select.Root>

    <Label>Value to send:</Label>
    <Input type="number" bind:value disabled={!editing} />

    <Label>Delay:</Label>
    <Input type="number" bind:value={minuteDelay} disabled={!editing} />
    <Label>minutes</Label>
</div>
