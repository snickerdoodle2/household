<script lang="ts">
    import * as Card from '$lib/components/ui/card';
    import { Label } from '$lib/components/ui/label';
    import * as Select from '$lib/components/ui/select';
    import FormInput from '@/components/FormInput.svelte';
    import { Button } from '@/components/ui/button';
    import { onMount } from 'svelte';
    import {
        newRuleSchema,
        ruleInternalSchema,
        type NewRule,
    } from '$lib/types/rule';
    import type { PageData } from './$types';
    import { authFetch } from '@/helpers/fetch';

    export let data: PageData;

    let loading = true;
    let sensors: { label: string; value: string }[] = [];
    let selectedSensor: { label: string; value: string };
    let rule: NewRule = {
        name: '',
        description: '',
        on_valid: {
            to: '',
            payload: {},
        },
        // @ts-expect-error nah dont wanna do this
        internal: {},
    };
    let errors: Record<string, string> = {};
    let internal = '';
    let payload = '';

    $: if (!loading && selectedSensor) {
        rule.on_valid.to = selectedSensor.value;
    }

    $: if (!loading) {
        try {
            rule.on_valid.payload = JSON.parse(payload);
            delete errors['payload'];
            errors = errors;
        } catch {
            errors['payload'] = 'Not a valid JSON';
        }
    }

    $: if (!loading) {
        try {
            const { data, success, error } = ruleInternalSchema.safeParse(
                JSON.parse(internal)
            );
            if (success) {
                delete errors['payload'];
                errors = errors;
                rule.internal = data;
            } else {
                console.log(error.issues);
            }
        } catch {
            errors['payload'] = 'Not a valid JSON';
        }
    }

    const handleSubmit = async () => {
        const { success, data, error } = newRuleSchema.safeParse(rule);

        if (!success) {
            console.log(error.issues);
            return;
        }

        const res = await authFetch(
            '/api/v1/rule',
            { method: 'POST', body: JSON.stringify(data) },
            fetch
        );

        console.log(await res.json());
    };

    onMount(async () => {
        sensors = (await data.sensors).map((e) => ({
            value: e.id,
            label: e.name,
        }));
        loading = false;
    });
</script>

{#if loading}
    <p>loading</p>
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
            />
            <FormInput
                name="description"
                type="text"
                label="Description"
                {errors}
                bind:value={rule.description}
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
            <Select.Root bind:selected={selectedSensor} required name="type">
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
            />
            <FormInput
                name="internal"
                type="text"
                label="Internal FIXME"
                {errors}
                bind:value={internal}
            />
        </Card.Content>
        <Card.Footer class="flex justify-end gap-3">
            <Button size="bold" on:click={handleSubmit}>Create</Button>
        </Card.Footer>
    </Card.Root>
{/if}
