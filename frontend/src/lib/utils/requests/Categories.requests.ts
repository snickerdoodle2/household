import { categoryStore } from '@/stores/Stores';
import type { Result } from '@/types/Result.types';
import { get } from 'svelte/store';
import { validateName } from '../Misc.utils';

export async function syncCategories() {
    categoryStore.set(['Kitchen', 'Living Room', 'Bedroom', 'Bathroom', 'Garage', 'Garden']); // temporary for now
}

export async function submitNewCategory(newCategory: string): Promise<Result<string, string>> {
    // Trim the new category
    const trimmedCategory = newCategory.trim();

    // Validate the new category
    const result = validateName(newCategory);
    if (result.isError) return result;

    // Get the current categories from the store
    const currentCategories = get(categoryStore);

    // Check if the category already exists
    const categoryExists = currentCategories.some(
        (category) => category.toLowerCase() === trimmedCategory.toLowerCase()
    );

    if (categoryExists) {
        return {
            isError: true,
            error: 'Category name already exists.',
        };
    }

    // TODO: send new category to the server / db

    // If the category is valid and does not exist, add it to the store
    categoryStore.update((categories) => [...categories, trimmedCategory]);

    // Return the new category as a success result
    return {
        isError: false,
        data: trimmedCategory,
    };
}
