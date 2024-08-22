import AddCategoryModal from '@/components/modals/categories/AddCategoryModal.svelte';
import ModifyCategoryModal from '@/components/modals/categories/ModifyCategoryModal.svelte';
import ConfirmationModal from '@/components/modals/ConfirmationModal.svelte';
import AddSensorModal from '@/components/modals/sensor/AddSensorModal.svelte';
import type { ComponentType } from 'svelte';

export enum ModalType {
    ADD_SENSOR = 'add_sensor',
    ADD_CATEGORY = 'add_category',
    MODIFY_CATEGORY = 'modify_category',
    CONFIRMATION_MODAL = 'confirmation_modal',
}
export type ModalDataPayload = {
    [ModalType.ADD_SENSOR]: undefined;
    [ModalType.ADD_CATEGORY]: undefined;
    [ModalType.MODIFY_CATEGORY]: string;
    [ModalType.CONFIRMATION_MODAL]: {
        message: string;
        acceptText: string;
        declineText: string;
        onAccept: () => void;
        onDecline: () => void;
    };
};

export type OpenedModalData<T extends ModalType> = {
    type: T;
    data: ModalDataPayload[T];
};

export const svelteModalMap: Record<ModalType, ComponentType> = {
    [ModalType.ADD_SENSOR]: AddSensorModal,
    [ModalType.ADD_CATEGORY]: AddCategoryModal,
    [ModalType.MODIFY_CATEGORY]: ModifyCategoryModal,
    [ModalType.CONFIRMATION_MODAL]: ConfirmationModal,
};

export function isModalData<T extends ModalType>(
    modalType: T,
    modalData: OpenedModalData<ModalType>
): modalData is OpenedModalData<T> {
    return modalData.type === modalType;
}
