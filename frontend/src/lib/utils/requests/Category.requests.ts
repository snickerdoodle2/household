import { categoryStore } from '@/stores/Stores';
import type { Result } from '@/types/Result.types';
import { validateCategory } from '../Category.utils';
import { get } from 'svelte/store';

export async function syncCategories() {
    categoryStore.set(['Kitchen', 'Living Room', 'Bedroom', 'Bathroom', 'Garage', 'Garden']); // temporary for now
}

export async function submitNewCategory(newCategory: string): Promise<Result<string, string>> {
    // Validate the new category
    const result = validateCategory(newCategory);
    if (result.isError) return result;

    // TODO: send new category to the server / db

    // If the category is valid and does not exist, add it to the store
    categoryStore.update((categories) => [...categories, newCategory]);

    // Return the new category as a success result
    return {
        isError: false,
        data: newCategory,
    };
}

export async function submitCategoryModification(
    oldName: string,
    newName: string
): Promise<Result<string, string>> {
    // Validate the new category
    const result = validateCategory(newName);
    if (result.isError) return result;

    // Check if the old category exists
    const categories = get(categoryStore);
    const categoryIndex = categories.indexOf(oldName);
    if (categoryIndex === -1) {
        return {
            isError: true,
            error: 'Modified category does not exist.',
        };
    }

    // Check if the new category already exists
    if (categories.includes(newName)) {
        return {
            isError: true,
            error: 'New category name already exists.',
        };
    }

    // Update the category in the store
    categoryStore.update((categories) => {
        categories[categoryIndex] = newName;
        return [...categories];
    });

    // Return the new category as a success result
    return {
        isError: false,
        data: newName,
    };
}

export async function submitCategoryDeletion(
    categoryToDelete: string
): Promise<Result<string, string>> {
    // Get the current categories from the store
    const categories = get(categoryStore);

    // Check if the category exists
    const categoryIndex = categories.indexOf(categoryToDelete);
    if (categoryIndex === -1) {
        return {
            isError: true,
            error: 'Category does not exist.',
        };
    }

    // Remove the category from the store
    categoryStore.update((categories) => {
        categories.splice(categoryIndex, 1);
        return [...categories];
    });

    // Return success result with the deleted category name
    return {
        isError: false,
        data: categoryToDelete,
    };
}
