<script lang="ts">
    import * as Select from '$lib/components/ui/select';
    import type { Selected } from 'bits-ui';
    import { onDestroy } from 'svelte';
    import { authToken } from '@/auth/token';
    import { get } from 'svelte/store';
    import type { LayoutData } from '../$types';

    let message = $state({});

    type Props = {
        data: LayoutData;
    };

    let { data }: Props = $props();

    let socket: WebSocket | undefined = undefined;

    let selected: string | undefined = $state();

    const updateSocket = (item: Selected<string> | undefined) => {
        const token = get(authToken);
        if (!token) return;
        if (!item || item.value.length === 0) return;
        if (socket) socket.close();

        message = {};
        selected = item.label;

        const url = new URL(
            `/api/v1/sensor/${item.value}/value`,
            window.location.href
        );
        url.protocol = url.protocol.replace('http', 'ws');
        url.searchParams.set('token', token.token);

        // TODO: auth not working with web socket
        socket = new WebSocket(url.toString());

        socket.addEventListener('message', (data) => {
            message = JSON.parse(data.data);
        });
    };

    onDestroy(() => {
        if (socket) socket.close();
    });
</script>

<div class="flex flex-col">
    <div>
        {#await data.sensors then sensors}
            <Select.Root
                items={sensors.map((e) => ({ value: e.id, label: e.name }))}
                onSelectedChange={updateSocket}
            >
                <Select.Trigger class="max-w-96">
                    <Select.Value placeholder="Select a sensor..." />
                </Select.Trigger>
                <Select.Content>
                    {#each sensors as sensor}
                        <Select.Item value={sensor.id} label={sensor.name}
                            >{sensor.name}</Select.Item
                        >
                    {/each}
                </Select.Content>
            </Select.Root>
            {#if selected}
                <p>Listening for sensor: <code>{selected}</code></p>
            {/if}
            {#if Object.keys(message).length > 0}}
                <code><pre>{JSON.stringify(message, null, 4)}</pre></code>
            {/if}
        {/await}
        <code><pre>{JSON.stringify($authToken, null, 4)}</pre></code>
        {#await data.user then user}
            <code><pre>{JSON.stringify(user, null, 4)}</pre></code>
        {/await}
    </div>
</div>
