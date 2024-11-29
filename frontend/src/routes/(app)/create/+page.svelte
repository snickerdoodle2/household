<script lang="ts">
    import { run, preventDefault } from 'svelte/legacy';
    import * as Card from '$lib/components/ui/card';
    import { Button } from '$lib/components/ui/button';
    import { Label } from '$lib/components/ui/label';
    import * as Select from '$lib/components/ui/select';
    import { newSensorSchema, sensorTypeSchema } from '$lib/types/sensor';
    import NewSensorInput from '$lib/components/FormInput.svelte';
    import { authFetch } from '@/helpers/fetch';
    import Input from '@/components/ui/input/input.svelte';
    import * as Dialog from '$lib/components/ui/dialog';
    import { goto, invalidate } from '$app/navigation';
    import { SENSOR_URL } from '@/helpers/sensor';

    const sensorTypes = sensorTypeSchema.options.map((e) => ({
        value: e,
        // TODO: Add capitalization or full lables
        label: e.replace('_', ' '),
    }));

    let name: string | undefined = $state(undefined);
    let refresh_rate: string = $state('');
    let uri: string | undefined = $state(undefined);
    let type: { value: string; label: string } | undefined = $state();
    let active: boolean = $state(false);
    let timeout: number;
    let errors: Partial<
        Record<'uri' | 'name' | 'refresh_rate' | 'type' | 'active', string>
    > = $state({});

    const getRefreshRate = () => {
        return active ? 1 : (refresh_rate ? +refresh_rate : undefined)
    }

    // eslint-disable-next-line @typescript-eslint/no-unsafe-function-type
    const debounce = (callback: Function, ...args: unknown[]) => {
        clearTimeout(timeout);
        timeout = window.setTimeout(() => callback(args), 300);
    };

    const validate = () => {
        const { success, error } = newSensorSchema.safeParse({
            name,
            refresh_rate: getRefreshRate(),
            uri,
            type: type?.value,
            active: typeof active === 'boolean' ? active : false,
        });
        
        
        if (!success) {
            errors = Object.fromEntries(
                error.issues.map((e) => [e.path[0], e.message])
            );
        } else {
            errors = {};
        }
        
        console.log('validate', error);

    };

    run(() => {
        debounce(validate, name, refresh_rate, uri, type, active);
    });

    const handleSubmit = async () => {
        const { success, data } = newSensorSchema.safeParse({
            name,
            refresh_rate: getRefreshRate(),
            uri,
            type: type?.value,
            active,
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
                        <Label
                            for={'refresh_rate'}
                            class="flex items-center justify-between text-base font-semibold"
                            >{'Refresh rate'}
                            {#if errors['refresh_rate'] && !active}
                                <span
                                    class="text-sm font-normal italic text-red-400"
                                    >{errors['refresh_rate']}</span
                                >
                            {/if}
                        </Label>

                        <div class="flex w-full flex-row items-center">
                            <div class="w-full">
                                <Input
                                    type="number"
                                    {name}
                                    bind:value={refresh_rate}
                                    required
                                    disabled={active}
                                    class={`${
                                        !active && errors['refresh_rate']
                                            ? 'border-2 border-red-600'
                                            : ''
                                    }`}
                                />
                            </div>

                            <div
                                class="ml-2 flex flex-row items-center justify-end"
                            >
                                <Label
                                    for="type"
                                    class="flex items-center justify-between text-base font-semibold"
                                >
                                    Active
                                    {#if errors['active']}
                                        <span
                                            class="text-sm font-normal italic text-red-400"
                                            >{errors['active']}</span
                                        >
                                    {/if}
                                </Label>
                                <Input
                                    type="checkbox"
                                    class="ml-2 w-8 {errors['active']
                                        ? 'border-2 border-red-600'
                                        : ''}"
                                    bind:checked={active}
                                />
                            </div>
                        </div>

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
