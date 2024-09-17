<script lang="ts">
    import * as Card from '$lib/components/ui/card';
    import { Label } from '$lib/components/ui/label';
    import { Input } from '$lib/components/ui/input';
    import { Button } from '$lib/components/ui/button';
    import * as Select from '$lib/components/ui/select';
    import { newSensorSchema, sensorTypeSchema } from '$lib/types/sensor';
    const labelClass = 'font-semibold text-base';

    const sensorTypes = sensorTypeSchema.options.map((e) => ({
        value: e,
        // TODO: Add capitalization or full lables
        label: e.replace('_', ' '),
    }));

    let name: string;
    let refresh_rate: number;
    let uri: string;
    let type: { value: string; label: string } | undefined;
    let timeout: number;
    let errors: Record<string, string> = {};

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
        errors = {};
    };

    $: debounce(validate, name, refresh_rate, uri, type);

    const handleSubmit = async () => {};
</script>

<pre>{JSON.stringify(errors, null, '    ')}</pre>
<Card.Root class="w-[512px] border-none shadow-none">
    <Card.Header class="text-3xl">
        <Card.Title>Create Sensor</Card.Title>
    </Card.Header>
    <form on:submit|preventDefault={handleSubmit}>
        <Card.Content class="grid grid-cols-[1fr_2fr] items-center gap-3">
            <Label for="name" class={labelClass}>Name</Label>
            <Input type="text" name="name" bind:value={name} required />
            <Label for="refresh_rate" class={labelClass}>Refresh Rate</Label>
            <Input
                type="number"
                name="refresh_rate"
                bind:value={refresh_rate}
                required
            />
            <Label for="uri" class={labelClass}>URI</Label>
            <Input type="text" name="uri" bind:value={uri} required />
            <Label for="sensor_type" class={labelClass}>Type</Label>
            <Select.Root bind:selected={type} required>
                <Select.Trigger>
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
                disabled={Object.keys(errors).length > 0}>Submit</Button
            >
        </Card.Footer>
    </form>
</Card.Root>
