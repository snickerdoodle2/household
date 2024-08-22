<script lang="ts">
    import { submitCategoryModification } from '@/utils/requests/Category.requests';
    import CategoryInputModal from './CategoryInputModal.svelte';
    import { get } from 'svelte/store';
    import { openedModalStore } from '@/stores/Stores';
    import { ModalType } from '@/types/Modal.types';
    import { closeModal, isModalData } from '@/utils/Modal.utils';
    let storeData = get(openedModalStore);
    let modifiedCategory = undefined as string | undefined;

    if (storeData && isModalData(ModalType.MODIFY_CATEGORY, storeData)) {
        modifiedCategory = storeData.data;
    } else {
        closeModal();
        modifiedCategory = '';
    }
</script>

<CategoryInputModal
    title="Modify Category"
    submit={(newCategoryName) => submitCategoryModification(modifiedCategory, newCategoryName)}
/>
