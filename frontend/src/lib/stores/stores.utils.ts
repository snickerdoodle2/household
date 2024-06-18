import type { ModalDataPayload, ModalType } from '@/types/modal';
import { openedModalStore } from './stores';

export function closeModal() {
    openedModalStore.set(null);
}

export function openModal<T extends ModalType>(
    type: T,
    data: ModalDataPayload[T]
) {
    console.log("Opening the modal", { type, data })
    openedModalStore.set({ type, data });
}
