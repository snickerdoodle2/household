<script lang="ts">
    import { run, preventDefault } from 'svelte/legacy';
    import * as Card from '$lib/components/ui/card';
    import { Button } from '$lib/components/ui/button';
    import { Label } from '$lib/components/ui/label';
    import * as Select from '$lib/components/ui/select';
    import { newSensorSchema, sensorTypeSchema } from '$lib/types/sensor';
    import NewSensorInput from '$lib/components/FormInput.svelte';
    import { authFetch } from '@/helpers/fetch';
    import * as Dialog from '$lib/components/ui/dialog';
    import { goto, invalidate } from '$app/navigation';
    import { SENSOR_URL } from '@/helpers/sensor';

    const sensorTypes = sensorTypeSchema.options.map((e) => ({
        value: e,
        // TODO: Add capitalization or full lables
        label: e.replace('_', ' '),
    }));

    let name: string = $state('');
    let refresh_rate: string = $state('');
    let uri: string = $state('');
    let type: { value: string; label: string } | undefined = $state();
    let timeout: number;
    let errors: Partial<
        Record<'uri' | 'name' | 'refresh_rate' | 'type', string>
    > = $state({});

    // eslint-disable-next-line @typescript-eslint/no-unsafe-function-type
    const debounce = (callback: Function, ...args: unknown[]) => {
        clearTimeout(timeout);
        timeout = window.setTimeout(() => callback(args), 300);
    };

    const validate = () => {
        const { success, error } = newSensorSchema.safeParse({
            name,
            refresh_rate: refresh_rate ? +refresh_rate : undefined,
            uri,
            type: type?.value,
        });

        if (!success) {
            errors = Object.fromEntries(
                error.issues.map((e) => [e.path[0], e.message])
            );
            return;
        }

        console.log('validate', errors);
        errors = {};
    };

    run(() => {
        debounce(validate, name, refresh_rate, uri, type);
    });

    const handleSubmit = async () => {
        const { success, data } = newSensorSchema.safeParse({
            name,
            refresh_rate: refresh_rate ? +refresh_rate : undefined,
            uri,
            type: type?.value,
        });

        if (!success) return;

        const res = await authFetch(
            '/api/v1/sensor',
            { method: 'POST', body: JSON.stringify(data) },
            fetch
        );

        const resJson = await res.json();
        console.log(resJson);

        if (!res.ok) {
            errors = resJson.error;
        } else {
            await invalidate(SENSOR_URL);
            close();
        }
    };

    const close = () => {
        goto(`/`);
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
            <Card.Root class="w-[600px] border-none shadow-none">
                <Card.Header class="text-3xl">
                    <Card.Title>Create Sensor</Card.Title>
                </Card.Header>
                <form onsubmit={preventDefault(handleSubmit)}>
                    <Card.Content
                        class="grid grid-cols-[3fr_4fr] items-center gap-3"
                    >
                        <NewSensorInput
                            name="name"
                            label="Name"
                            bind:value={name}
                            type="text"
                            {errors}
                        />
                        <NewSensorInput
                            name="refresh_rate"
                            label="Refresh rate"
                            bind:value={refresh_rate}
                            type="number"
                            {errors}
                        />
                        <NewSensorInput
                            name="uri"
                            label="URI"
                            bind:value={uri}
                            type="string"
                            {errors}
                        />
                        <Label
                            for="type"
                            class="flex items-center justify-between text-base font-semibold"
                        >
                            Type
                            {#if errors['type']}
                                <span
                                    class="text-sm font-normal italic text-red-400"
                                    >{errors['type']}</span
                                >
                            {/if}
                        </Label>
                        <Select.Root bind:selected={type} required name="type">
                            <Select.Trigger
                                class={errors['type']
                                    ? 'border-2 border-red-600'
                                    : ''}
                            >
                                <Select.Value />
                            </Select.Trigger>
                            <Select.Content>
                                {#each sensorTypes as type}
                                    <Select.Item value={type.value}
                                        >{type.label}</Select.Item
                                    >
                                {/each}
                            </Select.Content>
                        </Select.Root>
                    </Card.Content>
                    <Card.Footer class="flex justify-end gap-3">
                        <Button
                            size="bold"
                            type="submit"
                            disabled={Object.keys(errors).length > 0}
                            >Submit</Button
                        >
                    </Card.Footer>
                </form>
            </Card.Root>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>
