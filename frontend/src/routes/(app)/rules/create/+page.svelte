<script lang="ts">
    import { run } from 'svelte/legacy';
    import * as Dialog from '$lib/components/ui/dialog';
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
    import { goto, invalidate } from '$app/navigation';
    import RuleInternalBuilder from '@/components/rule/RuleInternalBuilder.svelte';
    import type { Sensor } from '@/types/sensor';
    import { RULE_URL } from '@/helpers/rule';
    import Input from '@/components/ui/input/input.svelte';

    type Props = {
        data: PageData;
    };

    let { data }: Props = $props();

    let loading = $state(true);
    let sensors: Sensor[] = $state([]);
    let selectedSensor: { label: string; value: string } = $state({
        label: '',
        value: '',
    });
    let rule: NewRule = $state({
        name: '',
        description: '',
        on_valid: {
            to: '',
            payload: {},
        },
        internal: {} as NewRule['internal'],
    });
    let errors: Record<string, string> = $state({});
    let internal = $state({});
    let payload: string = $state('');

    run(() => {
        if (!loading && selectedSensor) {
            rule.on_valid.to = selectedSensor.value;
        }
    });
    run(() => {
        if (!loading && payload) {
            rule.on_valid.payload = { value: Number(payload) };
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
        console.log(rule);
        const { success, data, error } = newRuleSchema.safeParse(rule);

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
            await invalidate(RULE_URL);
            close();
        }
    };

    onMount(async () => {
        sensors = await data.sensors;
        loading = false;
    });

    const close = () => {
        goto(`/rules/`);
    };
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
                <p>loading</p>
            {:else}
                <Card.Root class="min-w-[700px] border-none shadow-none">
                    <Card.Header class="text-3xl">
                        <Card.Title>New Rule</Card.Title>
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
                        />
                        <FormInput
                            name="description"
                            type="text"
                            label="Description"
                            {errors}
                            bind:value={rule.description}
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
                                step="0.1"
                                class={`w-full ${errors['on_valid.payload'] ? 'border-2 border-red-600' : ''}`}
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
                        />
                    </Card.Content>
                    <Card.Footer class="flex justify-end gap-3">
                        <Button size="bold" on:click={close}>Cancel</Button>
                        <Button size="bold" on:click={handleSubmit}
                            >Create</Button
                        >
                    </Card.Footer>
                </Card.Root>
            {/if}
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>
