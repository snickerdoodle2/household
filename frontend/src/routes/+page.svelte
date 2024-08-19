<script>
    import '../app.css';
    import { ModeWatcher } from 'mode-watcher';
    import MainPage from '@/components/MainPage.svelte';
    import { openedModalStore } from '@/stores/Stores';
    import { svelteModalMap } from '@/types/Modal.types';
    import AddSensorModal from '@/components/modals/AddSensorModal.svelte';
    import SensorInputModal from '@/components/modals/SensorInputModal.svelte';
    import { Modal } from 'flowbite-svelte';
</script>

<main>
    <div>
        <MainPage />
    </div>

    {#if $openedModalStore}
        <div class="absolute top-0 left-0 w-full h-full">
            <div class="flex items-center">
                <svelte:component
                    this={svelteModalMap[$openedModalStore.type]}
                />
            </div>
        </div>
    {/if}
    <Modal>
        <p>wtf</p>
        <p>wtf</p>
    </Modal>
</main>
