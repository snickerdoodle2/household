<script lang="ts">
import { SERVER_URL } from "$lib/const.js";
import { type Sensor } from "$lib/types/sensor";
let message = {};
let selectedId = "";

const WS_URL = SERVER_URL.replace("http", "ws");

// HACK: generating types not working :(
export let data: { ids: { data: Sensor[] } };

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
    <select bind:value={selectedId} placeholder="Select a sensor...">
        {#each data.ids.data as sensor}
            <option value={sensor.id}>{sensor.name}</option>
        {/each}
    </select>
    <p>Listening for sensor: <code>{selectedId}</code></p>
    <code>{JSON.stringify(message)}</code>
</main>

<style>
</style>
