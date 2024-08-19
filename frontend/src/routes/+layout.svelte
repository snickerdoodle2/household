<script>
    import '../app.css';
    import { ModeWatcher } from 'mode-watcher';
    import MainPage from '@/components/MainPage.svelte';
    import { openedModalStore } from '@/stores/Stores';
    import { svelteModalMap } from '@/types/Modal.types';
    import AddSensorModal from '@/components/modals/AddSensorModal.svelte';
</script>

<main>
    <ModeWatcher />
    <div>
        <MainPage />
    </div>
    {#if $openedModalStore}
        <div class="absolute w-full h-full">
            <div class="flex items-center">
                <svelte:component
                    this={svelteModalMap[$openedModalStore.type]}
                />
            </div>
        </div>
    {/if}
    <AddSensorModal />
</main>
