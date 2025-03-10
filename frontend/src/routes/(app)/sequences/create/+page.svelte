<script lang="ts">
    import * as Card from '$lib/components/ui/card';
    import FormInput from '@/components/FormInput.svelte';
    import { Button } from '@/components/ui/button';
    import { onMount } from 'svelte';
    import type { PageData } from './$types';
    import { authFetch } from '@/helpers/fetch';
    import { goto } from '$app/navigation';
    import {
        newSequenceSchema,
        type NewSequence,
        type SequenceAction,
    } from '@/types/sequence';
    import Label from '@/components/ui/label/label.svelte';
    import ActionsBuilder from '@/components/sequence/ActionsBuilder.svelte';
    import * as Dialog from '$lib/components/ui/dialog';

    type Props = {
        data: PageData;
    };
    let { data }: Props = $props();

    let loading = $state(true);
    let sensors: { label: string; value: string }[] = $state([]);
    let sequence: NewSequence = $state({
        name: '',
        description: '',
        actions: [] as SequenceAction[],
    });
    let errors: Record<string, string> = $state({});
    let actionFieldErrors: Record<number, string[]> = $state({});

    const handleSubmit = async () => {
        const { success, data, error } = newSequenceSchema.safeParse(sequence);

        if (!success) {
            error.issues.forEach((issue) => {
                const fieldPath = issue.path.join('.');
                if (fieldPath === 'name') {
                    errors['name'] = issue.message;
                } else if (fieldPath === 'description') {
                    errors['description'] = issue.message;
                }

                if (
                    issue.path[0] === 'actions' &&
                    typeof issue.path[1] == 'number' &&
                    typeof issue.path[2] == 'string'
                ) {
                    const errors = actionFieldErrors[issue.path[1]] ?? [];
                    actionFieldErrors[issue.path[1]] = [
                        ...errors,
                        issue.path[2],
                    ];
                }
            });
            console.log(error.issues);
            return;
        }

        const res = await authFetch(
            '/api/v1/sequence',
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
        sensors = (await data.sensors).map((e) => ({
            value: e.id,
            label: e.name,
        }));
        loading = false;
    });

    const leave = () => {
        goto(`/sequences/`);
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
                <Card.Root class="w-[1000px] border-none shadow-none">
                    <Card.Header class="text-3xl">
                        <Card.Title>New Sequence</Card.Title>
                    </Card.Header>
                    <Card.Content
                        class="grid grid-cols-[1fr_10fr] items-center gap-3"
                    >
                        <FormInput
                            name="name"
                            type="text"
                            label="Name"
                            {errors}
                            bind:value={sequence.name}
                        />
                        <FormInput
                            name="description"
                            type="text"
                            label="Description"
                            {errors}
                            bind:value={sequence.description}
                        />

                        <Label
                            for="type"
                            class="flex items-center justify-between text-base font-semibold"
                        >
                            Actions:
                        </Label>
                        <ActionsBuilder
                            bind:sensors
                            bind:actions={sequence.actions}
                            fieldErrors={actionFieldErrors}
                            editing={true}
                        />
                    </Card.Content>
                    <Card.Footer class="flex justify-end gap-3">
                        <Button size="bold" on:click={leave}>Cancel</Button>
                        <Button size="bold" on:click={handleSubmit}
                            >Create</Button
                        >
                    </Card.Footer>
                </Card.Root>
            {/if}
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>
