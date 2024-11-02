<script lang="ts">
    import * as Card from '$lib/components/ui/card';
    import FormInput from '@/components/FormInput.svelte';
    import type { PageData } from './$types';
    import { onMount } from 'svelte';
    import { Button } from '@/components/ui/button';
    import { authFetch } from '@/helpers/fetch';
    import { goto } from '$app/navigation';
    import { sequenceDetailsSchema, type SequenceAction, type SequenceDetails } from '@/types/sequence';
    import Label from '@/components/ui/label/label.svelte';
    import ActionsBuilder from '@/components/sequence/ActionsBuilder.svelte';
    export let data: PageData;
    let sequence: SequenceDetails;
    let loading = true;
    let errors: Record<string, string> = {};
    let editing = false;
    let sensors: { label: string; value: string }[] = [];
    let actions: SequenceAction[] = [];

    const leave = () => {
        goto(`/sequences/`);
    };

    const resetSequence = async () => {
        sequence = { ...(await data.sequence) };
        actions = JSON.parse(JSON.stringify(sequence.actions)); //deep copy
        console.log(actions)
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
            leave();
        }
    };

    const handleSubmit = async () => {
        const { data, success, error } = sequenceDetailsSchema.safeParse({
            ...sequence,
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
                    } else if (fieldPath === 'actions') {
                        errors['actions'] = issue.message;
                    }
                });
                console.log(error.issues);
                return;
            }
            return;
        }

        const { id, created_at, ...rest } = data; // eslint-disable-line @typescript-eslint/no-unused-vars
        console.log(rest);

        const res = await authFetch(`/api/v1/sequence/${sequence.id}`, {
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
        await resetSequence();
        loading = false;
    });
</script>

{#if loading}
    <p>Loading...</p>
{:else}
    <Card.Root class="w-[1000px] border-none shadow-none">
        <Card.Header class="text-3xl">
            <Card.Title>Sequence Details</Card.Title>
        </Card.Header>
        <Card.Content class="grid grid-cols-[1fr_5fr] items-center gap-3">
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
            <ActionsBuilder bind:sensors bind:actions={actions} bind:editing />
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
