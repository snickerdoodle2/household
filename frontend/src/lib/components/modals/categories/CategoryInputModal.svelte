<!-- src/routes/AddCategoryModal.svelte -->
<script lang="ts">
    import type { Result } from '@/types/Result.types';
    import { validateCategory } from '@/utils/Category.utils';
    import { closeModal } from '@/utils/Modal.utils';

    export let title;
    export let submit: (newCategory: string) => Promise<Result<string, string>>;
    export let category = '';
    let errorMessage = '';

    async function onSubmit() {
        // Clear previous error messages
        errorMessage = '';

        try {
            const result = await validateCategory(category);

            if (result.isError) {
                throw new Error(result.error);
            }
            submit(category);

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

<div class="card bg-background p-6 max-w-sm w-full relative">
    <button type="button" class="btn-dismiss absolute top-2 right-2" on:click={closeModal}>
        &times;
    </button>

    <p class="text-2xl font-bold mb-4">{title}</p>

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
            <button class="btn-a1" on:click={onSubmit}>Submit</button>
        </div>
    </div>
</div>
