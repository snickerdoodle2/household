<script>
    import '../app.css';
    import { ModeWatcher } from 'mode-watcher';
    import MainPage from '@/components/MainPage.svelte';
    import { openedModalStore } from '@/stores/Stores';
    import { svelteModalMap } from '@/types/Modal.types';
    import AddSensorModal from '@/components/modals/AddSensorModal.svelte';
    import SensorInputModal from '@/components/modals/SensorInputModal.svelte';
    import { Modal } from 'flowbite-svelte';
    import AddCategoryModal from '@/components/modals/AddCategoryModal.svelte';
</script>

<main>
    <div>
        <MainPage />
    </div>

    {#if $openedModalStore}
        <div class="container">
            <svelte:component this={svelteModalMap[$openedModalStore.type]} />
        </div>
    {/if}
</main>

<style>
    .container {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center; /* Vertically center the content */
        justify-content: center; /* Horizontally center the content */
        background-color: rgba(0, 0, 0, 0.7);
        transition: 0.2;
    }
</style>
