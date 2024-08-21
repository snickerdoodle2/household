import AddCategoryModal from '@/components/modals/AddCategoryModal.svelte';
import AddSensorModal from '@/components/modals/AddSensorModal.svelte';
import type { ComponentType } from 'svelte';

export enum ModalType {
    ADD_SENSOR = 'add_sensor',
    ADD_CATEGORY = 'add_category',
    // MONITOR_SENSOR = 'monitor',
}
export type ModalDataPayload = {
    [ModalType.ADD_SENSOR]: undefined;
    [ModalType.ADD_CATEGORY]: undefined;
    // [ModalType.MONITOR_SENSOR]: Sensor;
};

export type OpenedModalData<T extends ModalType> = {
    type: T;
    data: ModalDataPayload[T];
};

export const svelteModalMap: Record<ModalType, ComponentType> = {
    [ModalType.ADD_SENSOR]: AddSensorModal,
    [ModalType.ADD_CATEGORY]: AddCategoryModal
} 

