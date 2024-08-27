import { openedModalStore } from '@/stores/Stores';
import type { ModalDataPayload, ModalType, OpenedModalData } from '@/types/Modal.types';

export function openModal<T extends ModalType>(type: T, data: ModalDataPayload[T]) {
    console.log('Opening the modal', { type, data });
    openedModalStore.set({ type, data });
}

export function closeModal() {
    openedModalStore.set(null);
}

export function isModalData<T extends ModalType>(
    type: T,
    modalData: OpenedModalData<ModalType> | null
): modalData is OpenedModalData<T> {
    return modalData?.type === type;
}
