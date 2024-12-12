<script lang="ts">
    import * as Card from '$lib/components/ui/card';
    import FormInput from '@/components/FormInput.svelte';
    import type { PageData } from './$types';
    import { onMount } from 'svelte';
    import { Button } from '@/components/ui/button';
    import { authFetch } from '@/helpers/fetch';
    import { goto } from '$app/navigation';
    import {
        sequenceDetailsSchema,
        type SequenceAction,
        type SequenceDetails,
    } from '@/types/sequence';
    import Label from '@/components/ui/label/label.svelte';
    import ActionsBuilder from '@/components/sequence/ActionsBuilder.svelte';
    import * as Dialog from '$lib/components/ui/dialog';

    type Props = {
        data: PageData;
    };

    let { data }: Props = $props();

    let sequence: SequenceDetails = $state({
        id: '',
        name: '',
        description: '',
        actions: [],
        created_at: new Date(),
    });
    let loading = $state(true);
    let errors: Record<string, string> = $state({});
    let actionFieldErrors: Record<number, string[]> = $state({});
    let editing = $state(false);
    let sensors: { label: string; value: string }[] = $state([]);
    let actions: SequenceAction[] = $state([]);

    const close = () => {
        goto(`/sequences/`);
    };

    const resetSequence = async () => {
        sequence = { ...(await data.sequence) };
        actions = JSON.parse(JSON.stringify(sequence.actions)); //deep copy
        sequence.actions = actions;
        console.log(actions);
    };

    const handleCancel = async () => {
        await resetSequence();
        editing = false;
    };

    const handleDelete = async () => {
        const res = await authFetch(`/api/v1/sequence/${sequence.id}`, {
            method: 'DELETE',
        });

        console.log(await res.json());

        if (res.ok) {
            close();
        }
    };

    const handleStart = async () => {
        const res = await authFetch(`/api/v1/sequence/${sequence.id}/start`, {
            method: 'POST',
        });

        console.log(await res.json());

        if (res.ok) {
            close();
        }
    };

    const handleSubmit = async () => {
        const { data, success, error } = sequenceDetailsSchema.safeParse({
            ...sequence,
        });
        if (!success) {
            actionFieldErrors = {};
            errors = {};

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

        const { id, created_at, ...rest } = data; // eslint-disable-line @typescript-eslint/no-unused-vars
        console.log(rest);

        const res = await authFetch(`/api/v1/sequence/${sequence.id}`, {
            method: 'PUT',
            body: JSON.stringify(rest),
        });

        if (res.ok) {
            close();
        }

        console.log(await res.json());
    };

    onMount(async () => {
        sensors = (await data.sensors).map((e) => ({
            value: e.id,
            label: e.name,
        }));
        await resetSequence();
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
                <Card.Root class="w-[1000px] border-none shadow-none">
                    <Card.Header class="text-3xl">
                        <Card.Title>Sequence Details</Card.Title>
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
                            disabled={!editing}
                        />
                        <FormInput
                            name="description"
                            type="text"
                            label="Description"
                            {errors}
                            bind:value={sequence.description}
                            disabled={!editing}
                        />
                        <Label
                            for="type"
                            class="flex items-center justify-between text-base font-semibold"
                        >
                            Actions:
                        </Label>
                        <ActionsBuilder
                            bind:sensors
                            bind:actions
                            bind:editing
                            bind:fieldErrors={actionFieldErrors}
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
                            <Button on:click={handleStart} size="bold"
                                >Start</Button
                            >
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
