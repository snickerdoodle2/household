<script lang="ts">
    import { run } from 'svelte/legacy';
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
    import { goto, invalidate } from '$app/navigation';
    import RuleInternalBuilder from '@/components/rule/RuleInternalBuilder.svelte';
    import type { Sensor } from '@/types/sensor';
    import { RULE_URL } from '@/helpers/rule';
    import Input from '@/components/ui/input/input.svelte';
    import * as Dialog from '$lib/components/ui/dialog';

    type Props = {
        data: PageData;
    };

    let { data }: Props = $props();
    let rule: RuleDetails = $state({
        id: '',
        name: '',
        created_at: new Date(),
        description: '',
        on_valid: {
            to: '',
            payload: {},
        },
        internal: {} as RuleDetails['internal'],
    });
    let loading = $state(true);
    let errors: Record<string, string> = $state({});
    let editing = $state(false);
    let sensors: Sensor[] = $state([]);
    let selectedSensor: { label: string; value: string } = $state({
        label: '',
        value: '',
    });
    let internal = $state({});
    let payload = $state('');

    // TODO: make single validation function
    run(() => {
        if (!loading) {
            rule.on_valid.to = selectedSensor.value;
        }
    });

    run(() => {
        if (!loading) {
            try {
                const { data, success } =
                    ruleInternalSchema.safeParse(internal);
                if (success) {
                    rule.internal = data;
                }
            } catch {
                errors['internal'] = 'Invalid JSON';
            }
        }
    });

    run(() => {
        if (!loading) {
            try {
                rule.on_valid.payload = { value: Number(payload) };
                delete errors['payload'];
                errors = errors;
            } catch {
                errors['payload'] = 'Not a valid JSON';
            }
        }
    });

    const close = () => {
        goto(`/rules/`);
    };

    const resetRule = async () => {
        rule = { ...(await data.rule) };
        const sensor = sensors.find((e) => e.id === rule.on_valid.to);
        if (sensor) {
            selectedSensor = { value: sensor.id, label: sensor.name };
        }
        payload = JSON.stringify(rule.on_valid.payload['value']);
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
            await invalidate(RULE_URL);
            close();
        }
    };

    const handleSubmit = async () => {
        const { data, success, error } = ruleDetailsSchema.safeParse({
            ...rule,
        });
        if (!success) {
            if (!success) {
                error.issues.forEach((issue) => {
                    const fieldPath = issue.path.join('.');
                    if (fieldPath === 'name') {
                        errors['name'] = issue.message;
                    } else if (fieldPath === 'description') {
                        errors['description'] = issue.message;
                    } else if (fieldPath === 'on_valid.to') {
                        errors['on_valid.to'] = issue.message;
                    } else if (fieldPath === 'on_valid.payload') {
                        errors['on_valid.payload'] = issue.message;
                    } else if (fieldPath === 'internal') {
                        errors['internal'] = issue.message;
                    }
                });
                return;
            }
            return;
        }

        const { id, created_at, ...rest } = data; // eslint-disable-line @typescript-eslint/no-unused-vars

        console.log('Submitted rule modification:', rest, '(sent)');
        const res = await authFetch(`/api/v1/rule/${rule.id}`, {
            method: 'PUT',
            body: JSON.stringify(rest),
        });

        if (res.ok) {
            console.log('Rule modified:', await res.json(), '(recieved)');
            await invalidate(RULE_URL);
            close();
        } else {
            console.error(await res.json());
        }
    };

    onMount(async () => {
        sensors = await data.sensors;
        await resetRule();
        loading = false;
    });
</script>

<Dialog.Root
    open={true}
    onOpenChange={(opened) => {
        if (!opened) close();
    }}
>
    <Dialog.Portal>
        <Dialog.Overlay />
        <Dialog.Content
            class="flex max-w-none items-center justify-center px-8 py-4 md:w-fit"
        >
            {#if loading}
                <p>Loading...</p>
            {:else}
                <Card.Root class="min-w-[700px] border-none shadow-none">
                    <Card.Header class="text-3xl">
                        <Card.Title>Rule Details</Card.Title>
                    </Card.Header>
                    <Card.Content
                        class="grid grid-cols-[1fr_3fr] items-center gap-3"
                    >
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
                            name="on_valid.to"
                            class="pr-3 flex items-center justify-between text-base font-semibold"
                        >
                            Payload
                        </Label>
                        <div
                            class="flex grid-cols-3 gap-3 justify-center items-center"
                        >
                            <Input
                                type={'number'}
                                bind:value={payload}
                                required
                                class={`w-full ${errors['on_valid.payload'] ? 'border-2 border-red-600' : ''}`}
                                disabled={!editing}
                            />

                            <div class="flex items-center justify-center">
                                <Label
                                    for="type"
                                    class="flex items-center text-base font-semibold"
                                >
                                    to
                                    {#if errors['on_valid.to']}
                                        <span
                                            class="text-sm font-normal italic text-red-400"
                                            >{errors['type']}</span
                                        >
                                    {/if}
                                </Label>
                            </div>

                            <Select.Root
                                bind:selected={selectedSensor}
                                required
                                name="on_valid.to"
                                disabled={!editing}
                            >
                                <Select.Trigger
                                    class={`w-full ${errors['on_valid.to'] ? 'border-2 border-red-600' : ''}`}
                                >
                                    <Select.Value />
                                </Select.Trigger>
                                <Select.Content>
                                    {#each sensors as type}
                                        <Select.Item value={type.id}
                                            >{type.name}</Select.Item
                                        >
                                    {/each}
                                </Select.Content>
                            </Select.Root>
                        </div>
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
                            editingDisabled={!editing}
                        />
                    </Card.Content>
                    <Card.Footer class="flex justify-end gap-3">
                        {#if editing}
                            <Button
                                variant="destructive"
                                size="bold"
                                on:click={handleDelete}>Delete</Button
                            >
                            <Button
                                variant="outline"
                                size="bold"
                                on:click={handleCancel}>Cancel</Button
                            >
                            <Button size="bold" on:click={handleSubmit}
                                >Submit</Button
                            >
                        {:else}
                            <Button on:click={close} size="bold">Cancel</Button>
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
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>
