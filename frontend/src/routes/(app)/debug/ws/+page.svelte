<script lang="ts">
    import { SensorWebsocket } from '@/helpers/socket.svelte';
    let ws = new SensorWebsocket();
</script>

<p>Websocket</p>
<p>{JSON.stringify(ws.data)}</p>
