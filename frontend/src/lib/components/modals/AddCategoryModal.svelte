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

<main>
    <p class="text-primary-foreground text-2xl font-bold justify-center">Add new category</p>
    <div class="card flex flex-col p-4">
        <div class="flex items-center">
            <input
                type="text"
                class="input-field mr-2"
                placeholder="Enter your category name here"
                bind:value={category}
            />
            <button class="btn-primary" on:click={submit}>Submit</button>
        </div>
        {#if errorMessage}
            <p class="text-error mt-2">{errorMessage}</p>
        {/if}
    </div>
</main>
