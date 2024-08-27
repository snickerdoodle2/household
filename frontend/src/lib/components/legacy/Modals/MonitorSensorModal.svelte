<!-- src/routes/AddSensorForm.svelte -->
<script lang="ts">
    import { Modal } from 'flowbite-svelte';
    import { getWSUrl } from '@/config/const';
    import { onDestroy } from 'svelte';
    import { openedModalStore } from '@/stores/Stores';
    import { ModalType, isModalData } from '@/types/Modal.types';
    import type { Sensor } from '@/types/Sensor.types';

    let socket: WebSocket | undefined = undefined;

    let open = false;

    openedModalStore.subscribe((value) => {
        if (isModalData(ModalType.MONITOR_SENSOR, value)) {
            open = true;
            updateSocket(value.data);
        } else {
            open = false;
            socket = undefined;
        }
    });

    const WS_URL = getWSUrl();

    let message = {};

    const updateSocket = (sensor: Sensor) => {
        if (socket) socket.close();

        message = {};

        socket = new WebSocket(`${WS_URL}/api/v1/sensor/${sensor.id}/value`);

        socket.addEventListener('message', (data) => {
            message = JSON.parse(data.data);
        });
    };

    onDestroy(() => {
        if (socket) socket.close();
    });
</script>

<Modal title={'Monitoring Sensor'} autoclose={false} bind:open>
    <code><pre>{JSON.stringify(message, null, 4)}</pre></code>
</Modal>

<style>
</style>
