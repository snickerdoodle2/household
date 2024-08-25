<script>
    import '../app.css';
    import MainPage from '@/components/MainPage.svelte';
    import { openedModalStore } from '@/stores/Stores';
    import { svelteModalMap } from '@/types/Modal.types';
</script>

<main>
    <MainPage />
    {#if $openedModalStore}
        <div class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
            <svelte:component this={svelteModalMap[$openedModalStore.type]} />
        </div>
    {/if}
</main>
