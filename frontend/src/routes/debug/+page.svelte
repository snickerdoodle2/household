<script lang="ts">
import * as Select from "$lib/components/ui/select";
import { getWSUrl } from "@/const";
import type { Selected } from "bits-ui";
import { onDestroy } from "svelte";
import type { PageData } from "./$types";
import { authToken } from "@/auth/token";

const WS_URL = getWSUrl();

let message = {};

export let data: PageData;

const { user } = data;

let socket: WebSocket | undefined = undefined;

let selected: string | undefined;

const updateSocket = (item: Selected<string> | undefined) => {
    if (!item || item.value.length === 0) return;
    if (socket) socket.close();

    message = {};
    selected = item.label;

    // TODO: auth not working with web socket
    socket = new WebSocket(`${WS_URL}/api/v1/sensor/${item.value}/value`);

    socket.addEventListener("message", (data) => {
        message = JSON.parse(data.data);
    });
};

onDestroy(() => {
    if (socket) socket.close();
});
</script>

<div class="flex flex-col">
    <div>
    <Select.Root items={data.sensors} onSelectedChange={updateSocket}>
        <Select.Trigger class="max-w-96">
            <Select.Value placeholder="Select a sensor..." />
        </Select.Trigger>
        <Select.Content>
            {#each data.sensors as sensor}
                <Select.Item value={sensor.value} label={sensor.label}
                    >{sensor.label}</Select.Item
                >
            {/each}
        </Select.Content>
    </Select.Root>
    {#if selected}
        <p>Listening for sensor: <code>{selected}</code></p>
    {/if}
    <code><pre>{JSON.stringify(message, null, 4)}</pre></code>
    {#if $authToken}
    <code><pre>{JSON.stringify($authToken, null, 4)}</pre></code>
    {/if}
    {#if user}
    <code><pre>{JSON.stringify(user, null, 4)}</pre></code>
    {/if}
    </div>
</div>
