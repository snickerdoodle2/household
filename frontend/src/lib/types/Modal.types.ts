import AddSensorModal from '@/components/modals/AddSensorModal.svelte';
import type { ComponentType } from 'svelte';

export enum ModalType {
    ADD_SENSOR = 'add',
    // MODIFY_SENSOR = 'modify',
    // MONITOR_SENSOR = 'monitor',
}
export type ModalDataPayload = {
    [ModalType.ADD_SENSOR]: undefined;
    // [ModalType.MODIFY_SENSOR]: Sensor;
    // [ModalType.MONITOR_SENSOR]: Sensor;
};

export type OpenedModalData<T extends ModalType> = {
    type: T;
    data: ModalDataPayload[T];
};

export const svelteModalMap: Record<ModalType, ComponentType> = {
    [ModalType.ADD_SENSOR]: AddSensorModal
} 

