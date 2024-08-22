<!-- src/routes/WarningModal.svelte -->
<script lang="ts">
    import { openedModalStore } from '@/stores/Stores';
    import { ModalType } from '@/types/Modal.types';
    import { closeModal, isModalData } from '@/utils/Modal.utils';
    import { get } from 'svelte/store';

    function getStoreData() {
        let storeData = get(openedModalStore);
        if (storeData && isModalData(ModalType.CONFIRMATION_MODAL, storeData)) {
            return storeData.data;
        } else {
            closeModal();
            return {
                message: 'Are you sure you want to proceed?',
                acceptText: 'Accept',
                declineText: 'Decline',
                onAccept: () => {},
                onDecline: () => {},
            };
        }
    }

    let { message, acceptText, declineText, onAccept, onDecline } = getStoreData();

    function handleAccept() {
        onAccept();
        closeModal();
    }

    function handleDecline() {
        onDecline();
        closeModal();
    }
</script>

<main class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
    <div class="card bg-background p-6 max-w-sm w-full relative">
        <!-- Close button -->
        <button type="button" class="btn-exit absolute top-2 right-2" on:click={handleDecline}>
            &times;
        </button>

        <p class="text-xl font-bold mb-4">Warning</p>

        <div class="flex flex-col space-y-4">
            <p class="card text-primary-foreground">{message}</p>
            <div class="flex justify-end gap-2">
                <button class="btn-primary" on:click={handleAccept}>{acceptText}</button>
                <button class="btn-secondary" on:click={handleDecline}>{declineText}</button>
            </div>
        </div>
    </div>
</main>

<style>
    /* Example styles - adjust based on your project's design system */
    .btn-primary {
        background-color: hsl(var(--primary));
        color: hsl(var(--primary-foreground));
        padding: 0.5rem 1rem;
        border-radius: var(--radius);
    }

    .btn-secondary {
        background-color: hsl(var(--secondary));
        color: hsl(var(--secondary-foreground));
        padding: 0.5rem 1rem;
        border-radius: var(--radius);
    }
</style>
