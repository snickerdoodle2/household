<script lang="ts">
    import * as Card from '$lib/components/ui/card';
    import { Label } from '$lib/components/ui/label';
    import * as Select from '$lib/components/ui/select';
    import FormInput from '@/components/FormInput.svelte';
    import {
        type RuleDetails,
        ruleDetailsSchema,
        ruleInternalSchema,
    } from '@/types/rule';
    import type { PageData } from './$types';
    import { onMount } from 'svelte';
    import { Button } from '@/components/ui/button';
    import { authFetch } from '@/helpers/fetch';
    import { goto } from '$app/navigation';
    export let data: PageData;
    let rule: RuleDetails;
    let loading = true;
    let errors: Record<string, string> = {};
    let editing = false;
    let sensors: { label: string; value: string }[] = [];
    let selectedSensor: { label: string; value: string };
    let internal = '';
    let payload = '';

    // TODO: make single validation function
    $: if (!loading) {
        rule.on_valid.to = selectedSensor.value;
    }

    $: if (!loading) {
        try {
            const { data, success } = ruleInternalSchema.safeParse(
                JSON.parse(internal)
            );
            if (success) {
                rule.internal = data;
            }
        } catch {
            errors['internal'] = 'Invalid JSON';
        }
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

    const leave = () => {
        goto(`/rules/`);
    };

    const resetRule = async () => {
        rule = { ...(await data.rule) };
        const tmp = sensors.find((e) => e.value === rule.on_valid.to);
        if (tmp) {
            selectedSensor = tmp;
        }
        payload = JSON.stringify(rule.on_valid.payload);
        internal = JSON.stringify(rule.internal);
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

        if (res.ok) {
            leave();
        }
    };

    const handleSubmit = async () => {
        const { data, success, error } = ruleDetailsSchema.safeParse({
            ...rule,
        });
        if (!success) {
            console.log(error.issues);
            if (!success) {
                error.issues.forEach((issue) => {
                    const fieldPath = issue.path.join('.');
                    if (fieldPath === 'name') {
                        errors['name'] = issue.message;
                    } else if (fieldPath === 'description') {
                        errors['description'] = issue.message;
                    } else if (fieldPath === 'on_valid.to') {
                        errors['type'] = issue.message;
                    } else if (fieldPath === 'on_valid.payload') {
                        errors['payload'] = issue.message;
                    } else if (fieldPath === 'internal') {
                        errors['internal'] = issue.message;
                    }
                });
                console.log(error.issues);
                return;
            }
            return;
        }

        const { id, created_at, ...rest } = data; // eslint-disable-line @typescript-eslint/no-unused-vars
        console.log(rest);

        const res = await authFetch(`/api/v1/rule/${rule.id}`, {
            method: 'PUT',
            body: JSON.stringify(rest),
        });

        if (res.ok) {
            leave();
        }

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
                <Button size="bold" on:click={handleSubmit}>Submit</Button>
            {:else}
                <Button on:click={leave} size="bold">Cancel</Button>
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
