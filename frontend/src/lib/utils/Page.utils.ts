import { currentPageStore } from '@/stores/Stores';
import type { PageType } from '@/types/Page.types';

export function openPage(page: PageType) {
    currentPageStore.set(page);
}
