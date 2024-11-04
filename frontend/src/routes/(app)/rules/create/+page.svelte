<script lang="ts">
    import { run } from 'svelte/legacy';

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
    import { goto } from '$app/navigation';
    import RuleInternalBuilder from '@/components/rule/RuleInternalBuilder.svelte';
    import type { Sensor } from '@/types/sensor';

    type Props = {
        data: PageData;
    };

    let { data }: Props = $props();

    let loading = $state(true);
    let sensors: Sensor[] = $state([]);
    let selectedSensor: { label: string; value: string } = $state();
    let rule: NewRule = $state({
        name: '',
        description: '',
        on_valid: {
            to: '',
            payload: {},
        },
        // @ts-expect-error nah dont wanna do this
        internal: {},
    });
    let errors: Record<string, string> = $state({});
    let internal = {};
    let payload = $state('');

    run(() => {
        if (!loading && selectedSensor) {
            rule.on_valid.to = selectedSensor.value;
        }
    });

    run(() => {
        if (!loading) {
            try {
                rule.on_valid.payload = JSON.parse(payload);
                delete errors['payload'];
                errors = errors;
            } catch {
                errors['payload'] = 'Not a valid JSON';
            }
        }
    });

    run(() => {
        if (!loading) {
            try {
                const { data, success, error } =
                    ruleInternalSchema.safeParse(internal);

                if (success) {
                    delete errors['internal'];
                    errors = errors;
                    rule.internal = data;
                } else {
                    console.log('tbhjasr', error.issues);
                }
            } catch {
                errors['internal'] = 'Not a valid JSON';
            }
        }
    });

    const handleSubmit = async () => {
        const { success, data, error } = newRuleSchema.safeParse(rule);

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

        const res = await authFetch(
            '/api/v1/rule',
            { method: 'POST', body: JSON.stringify(data) },
            fetch
        );

        console.log(await res.json());
        if (!res.ok) {
            // TODO: direct errors to proper fields
            console.log('error');
        } else {
            leave();
        }
    };

    onMount(async () => {
        sensors = await data.sensors;
        loading = false;
    });

    const leave = () => {
        goto(`/rules/`);
    };
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
                        <Select.Item value={type.id}>{type.name}</Select.Item>
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

            <Label
                for="type"
                class="flex items-center justify-between text-base font-semibold"
            >
                Internal:
            </Label>
            <RuleInternalBuilder
                bind:internal={rule.internal}
                {sensors}
                bind:parent={rule}
                secondParent={undefined}
            />
        </Card.Content>
        <Card.Footer class="flex justify-end gap-3">
            <Button size="bold" on:click={leave}>Cancel</Button>
            <Button size="bold" on:click={handleSubmit}>Create</Button>
        </Card.Footer>
    </Card.Root>
{/if}
