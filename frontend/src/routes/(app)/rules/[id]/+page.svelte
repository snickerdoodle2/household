<script lang="ts">
    import { run } from 'svelte/legacy';
    import * as Card from '$lib/components/ui/card';
    import { Label } from '$lib/components/ui/label';
    import * as Select from '$lib/components/ui/select';
    import FormInput from '@/components/FormInput.svelte';
    import { type RuleDetails, ruleDetailsSchema } from '@/types/rule';
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
    import type { Sequence } from '@/types/sequence';

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
            target_type: 'sensor',
            target_id: '',
            payload: {},
        },
        internal: {} as RuleDetails['internal'],
    });
    let loading = $state(true);
    let errors: Record<string, string> = $state({});
    let editing = $state(false);
    let sensors: Sensor[] = $state([]);
    let sequences: Sequence[] = $state([]);
    let selectedSensor: { label: string; value: string } = $state({
        label: '',
        value: '',
    });
    let payload = $state('');
    let selectedSequence: { label: string; value: string } = $state({
        label: '',
        value: '',
    });
    let isSensorPayload: boolean = $state(true);

    // TODO: make single validation function
    run(() => {
        if (loading) return;

        if (isSensorPayload) {
            if (payload && selectedSensor.value !== '') {
                rule.on_valid = {
                    target_type: 'sensor',
                    target_id: selectedSensor.value,
                    payload: { value: Number(payload) },
                };
            }
        } else {
            if (selectedSequence.value !== '') {
                rule.on_valid = {
                    target_type: 'sequence',
                    target_id: selectedSequence.value,
                    payload: {},
                };
            }
        }
    });

    const close = () => {
        goto(`/rules/`);
    };

    const resetRule = async () => {
        try {
            rule = { ...(await data.rule) } as RuleDetails;

            loading = false;
            errors = {};

            if (rule.on_valid.target_type === 'sensor') {
                const sensor = sensors.find(
                    (e) => e.id === rule.on_valid.target_id
                );

                if (sensor) {
                    selectedSensor = { value: sensor.id, label: sensor.name };
                    payload = JSON.stringify(rule.on_valid.payload['value']);
                } else {
                    selectedSensor = { value: '', label: '' };
                    payload = '';
                }
                isSensorPayload = true;
            } else if (rule.on_valid.target_type === 'sequence') {
                const sequence = sequences.find(
                    (e) => e.id === rule.on_valid.target_id
                );

                if (sequence) {
                    selectedSequence = {
                        value: sequence.id,
                        label: sequence.name,
                    };
                } else {
                    selectedSequence = { value: '', label: '' };
                }
                payload = '';
                isSensorPayload = false;
            }
        } catch (error) {
            errors = { fetch: 'Failed to fetch rule data.' };
            console.error('Error in resetRule:', error);
        } finally {
            loading = false;
        }
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
                    } else if (fieldPath === 'on_valid.target_id') {
                        errors['on_valid.target_id'] = issue.message;
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
        sequences = await data.sequences;
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
                            name="on_valid.target_id"
                            class="pr-3 flex items-center justify-between text-base font-semibold"
                        >
                            Payload
                        </Label>
                        <div
                            class="flex grid-cols-3 gap-3 justify-center items-center"
                        >
                            <div class="flex">
                                <Button
                                    disabled={isSensorPayload || !editing}
                                    on:click={() => {
                                        isSensorPayload = true;
                                    }}>Sensor</Button
                                >
                                <Button
                                    disabled={!isSensorPayload || !editing}
                                    on:click={() => {
                                        isSensorPayload = false;
                                    }}>Sequence</Button
                                >
                            </div>

                            {#if isSensorPayload}
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
                                        {#if errors['on_valid.target_id']}
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
                                    name="on_valid.target_id"
                                    disabled={!editing}
                                >
                                    <Select.Trigger
                                        class={`w-full ${errors['on_valid.target_id'] ? 'border-2 border-red-600' : ''}`}
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
                            {:else}
                                <Select.Root
                                    bind:selected={selectedSequence}
                                    required
                                    name="on_valid.target_id"
                                    disabled={!editing}
                                >
                                    <Select.Trigger
                                        class={`w-full ${errors['on_valid.target_id'] ? 'border-2 border-red-600' : ''}`}
                                    >
                                        <Select.Value />
                                    </Select.Trigger>
                                    <Select.Content>
                                        {#each sequences as sequence}
                                            <Select.Item value={sequence.id}
                                                >{sequence.name}</Select.Item
                                            >
                                        {/each}
                                    </Select.Content>
                                </Select.Root>
                            {/if}
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
                            {#if data.currentUser.role === 'admin'}
                                <Button
                                    on:click={() => {
                                        editing = true;
                                    }}
                                    size="bold">Edit</Button
                                >
                            {/if}
                        {/if}
                    </Card.Footer>
                </Card.Root>
            {/if}
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>
