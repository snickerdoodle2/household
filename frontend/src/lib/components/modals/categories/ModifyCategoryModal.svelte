<script lang="ts">
    import { submitCategoryModification } from '@/utils/requests/Category.requests';
    import CategoryInputModal from './CategoryInputModal.svelte';
    import { get } from 'svelte/store';
    import { openedModalStore } from '@/stores/Stores';
    import { ModalType } from '@/types/Modal.types';
    import { closeModal, isModalData } from '@/utils/Modal.utils';
    let storeData = get(openedModalStore);
    let modifiedCategory = getModifiedCategory();

    function getModifiedCategory() {
        if (storeData && isModalData(ModalType.MODIFY_CATEGORY, storeData)) {
            return storeData.data;
        } else {
            closeModal();
            return '';
        }
    }
</script>

<CategoryInputModal
    title="Modify Category"
    category={modifiedCategory}
    submit={(newCategoryName) => submitCategoryModification(modifiedCategory, newCategoryName)}
/>
