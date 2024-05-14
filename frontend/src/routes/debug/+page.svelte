<script lang="ts">
    import { Input } from "$lib/components/ui/input/index.js";
    let message = {};
    let id = "";

    let socket: WebSocket | undefined = undefined;

    const updateSocket = (id: string) => {
        if (id.length === 0) return;
        if (socket) socket.close();

        socket = new WebSocket(`ws://localhost:8080/api/v1/sensor/${id}/value`);

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
        updateSocket(id);
    }
</script>

<main class="w-screen h-screen flex flex-col justify-center items-center">
    <Input class="max-w-xs" bind:value={id} />

    <p>Listening for sensor: <code>{id}</code></p>
    <code>{JSON.stringify(message)}</code>
</main>

<style>
</style>
