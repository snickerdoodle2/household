<script lang="ts">
    import * as Card from '$lib/components/ui/card';
    import { Label } from '$lib/components/ui/label';
    import * as Select from '$lib/components/ui/select';
    import FormInput from '@/components/FormInput.svelte';
    import type { RuleDetails } from '@/types/rule';
    import type { PageData } from './$types';
    import { onMount } from 'svelte';
    import { Button } from '@/components/ui/button';
    import { authFetch } from '@/helpers/fetch';
    export let data: PageData;
    let rule: RuleDetails;
    let loading = true;
    let errors: Record<string, string> = {};
    let editing = false;
    let sensors: { label: string; value: string }[] = [];
    let selectedSensor: { label: string; value: string };
    let internal = '';
    let payload = '';

    $: if (!loading) {
        rule.on_valid.to = selectedSensor.value;
    }
    const resetRule = async () => {
        rule = { ...(await data.rule) };
        const tmp = sensors.find((e) => e.value === rule.on_valid.to);
        if (tmp) {
            selectedSensor = tmp;
        }
        payload = JSON.stringify(rule.on_valid.payload);
        internal = JSON.stringify(rule);
    };

    const handleCancel = async () => {
        await resetRule();
        editing = false;
    };

    const handleDelete = async () => {
        // TODO: ask for confirmation!!!
        const res = await authFetch(`/api/v1/rule/${rule.id}`, {
            method: 'DELETE',
        });
        console.log(await res.json());
    };

    onMount(async () => {
        sensors = (await data.sensors).map((e) => ({
            value: e.id,
            label: e.name,
        }));
        await resetRule();
        loading = false;
    });
</script>

{#if loading}
    <p>Loading...</p>
{:else}
    <Card.Root class="w-[600px] border-none shadow-none">
        <Card.Header class="text-3xl">
            <Card.Title>Rule Details</Card.Title>
        </Card.Header>
        <Card.Content class="grid grid-cols-[1fr_2fr] items-center gap-3">
            <FormInput
                name="name"
                type="text"
                label="Name"
                {errors}
                bind:value={rule.name}
                disabled={!editing}
            />
            <FormInput
                name="description"
                type="text"
                label="Description"
                {errors}
                bind:value={rule.description}
                disabled={!editing}
            />
            <Label
                for="type"
                class="flex items-center justify-between text-base font-semibold"
            >
                To
                {#if errors['type']}
                    <span class="text-sm font-normal italic text-red-400"
                        >{errors['type']}</span
                    >
                {/if}
            </Label>
            <Select.Root
                bind:selected={selectedSensor}
                required
                name="type"
                disabled={!editing}
            >
                <Select.Trigger
                    class={errors['type'] ? 'border-2 border-red-600' : ''}
                >
                    <Select.Value />
                </Select.Trigger>
                <Select.Content>
                    {#each sensors as type}
                        <Select.Item value={type.value}
                            >{type.label}</Select.Item
                        >
                    {/each}
                </Select.Content>
            </Select.Root>
            <FormInput
                name="payload"
                type="text"
                label="Payload"
                {errors}
                bind:value={payload}
                disabled={!editing}
            />
            <FormInput
                name="internal"
                type="text"
                label="Internal FIXME"
                {errors}
                bind:value={internal}
                disabled={!editing}
            />
        </Card.Content>
        <Card.Footer class="flex justify-end gap-3">
            {#if editing}
                <Button
                    variant="destructive"
                    size="bold"
                    on:click={handleDelete}>Delete</Button
                >
                <Button variant="outline" size="bold" on:click={handleCancel}
                    >Cancel</Button
                >
                <Button size="bold">Submit</Button>
            {:else}
                <Button
                    on:click={() => {
                        editing = true;
                    }}
                    size="bold">Edit</Button
                >
            {/if}
        </Card.Footer>
    </Card.Root>
{/if}
