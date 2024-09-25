<script lang="ts">
    import * as Card from '$lib/components/ui/card';
    import FormInput from '@/components/FormInput.svelte';
    import type { RuleDetails } from '@/types/rule';
    import type { PageData } from './$types';
    import { onMount } from 'svelte';
    import { Button } from '@/components/ui/button';
    export let data: PageData;
    let rule: RuleDetails;
    let loading = true;
    let errors: Record<string, string> = {};
    let editing = false;
    let internal = '';

    const handleCancel = async () => {
        rule = { ...(await data.rule) };
        internal = JSON.stringify(rule);
        editing = false;
    };

    onMount(async () => {
        rule = { ...(await data.rule) };
        internal = JSON.stringify(rule);
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
                <Button variant="destructive" size="bold">Delete</Button>
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
