<!-- src/routes/AddCategoryModal.svelte -->
<script lang="ts">
    import { closeModal } from '@/utils/Modal.utils';
    import { submitNewCategory } from '@/utils/requests/Categories.requests';

    let category = '';
    let errorMessage = '';

    async function submit() {
        // Clear previous error messages
        errorMessage = '';

        try {
            const result = await submitNewCategory(category);

            if (result.isError) {
                throw new Error(result.error);
            }

            // Clear the input field or redirect
            category = '';
            closeModal();
        } catch (error) {
            // Type assertion to Error
            if (error instanceof Error) {
                // Handle known Error types
                console.error('Error:', error.message);
                errorMessage = error.message;
            } else {
                // Handle unknown error types
                console.error('An unknown error occurred:', error);
                errorMessage = 'An unknown error occurred.';
            }
        }
    }
</script>

<main class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
    <div class="card bg-background p-6 max-w-sm w-full relative">
        <!-- Close button -->
        <button type="button" class="btn-exit absolute top-2 right-2" on:click={closeModal}>
            &times;
        </button>

        <p class="text-2xl font-bold mb-4">Add new category</p>

        <div class="flex flex-col space-y-4">
            <input
                type="text"
                class="input-field w-full"
                placeholder="Enter your category name here"
                bind:value={category}
            />
            {#if errorMessage}
                <p class="text-error mt-2">{errorMessage}</p>
            {/if}
            <div class="flex justify-end gap-2">
                <button class="btn-primary" on:click={submit}>Submit</button>
            </div>
        </div>
    </div>
</main>
