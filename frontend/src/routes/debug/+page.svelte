<script lang="ts">
import { SERVER_URL } from "$lib/const.js";
import * as Select from "$lib/components/ui/select";

import type { PageData } from "./$types";

let message = {};
let selectedId = "";

const WS_URL = SERVER_URL.replace("http", "ws");

export let data: PageData;

let socket: WebSocket | undefined = undefined;

const updateSocket = (id: string) => {
	if (id.length === 0) return;
	if (socket) socket.close();

	message = {};

	socket = new WebSocket(`${WS_URL}/api/v1/sensor/${id}/value`);

	socket.addEventListener("open", () => {
		console.log("Opened");
	});

	socket.addEventListener("message", (data) => {
		message = JSON.parse(data.data);
	});

	socket.addEventListener("close", () => {
		console.log("Closed");
	});
};

$: {
	updateSocket(selectedId);
}
</script>

<main class="w-screen h-screen flex flex-col justify-center items-center">
    <Select.Root>
        <Select.Trigger>
            <Select.Value placeholder="Select a sensor..." />
        </Select.Trigger>
        <Select.Content>
            {#each data.ids.data as sensor}
                <Select.Item value={sensor.id}>{sensor.name}</Select.Item>
            {/each}
        </Select.Content>
    </Select.Root>
    <p>Listening for sensor: <code>{selectedId}</code></p>
    <code>{JSON.stringify(message)}</code>
</main>

<style>
</style>
