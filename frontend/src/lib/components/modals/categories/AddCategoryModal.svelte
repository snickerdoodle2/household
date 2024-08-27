<script>
    import { submitNewCategory } from '@/utils/requests/Category.requests';
    import CategoryInputModal from './CategoryInputModal.svelte';
</script>

<CategoryInputModal title="Add Category" submit={submitNewCategory} />
