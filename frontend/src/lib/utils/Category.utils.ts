import { get } from 'svelte/store';
import { validateName } from './Misc.utils';
import { categoryStore } from '@/stores/Stores';
import type { Result } from '@/types/Result.types';

export function validateCategory(newCategory: string): Result<string, string> {
    const result = validateName(newCategory);
    if (result.isError) return result;

    // Get the current categories from the store
    const currentCategories = get(categoryStore);

    // Check if the category already exists
    const categoryExists = currentCategories.some(
        (category) => category.toLowerCase() === newCategory.toLowerCase()
    );

    if (categoryExists) {
        return {
            isError: true,
            error: 'Category name already exists.',
        };
    }

    return {
        isError: false,
        data: newCategory,
    };
}
