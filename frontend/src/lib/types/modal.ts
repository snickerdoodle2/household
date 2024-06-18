import type { Sensor } from './sensor';

export enum ModalType {
    ADD_SENSOR = 'add',
    MODIFY_SENSOR = 'modify',
    MONITOR_SENSOR = 'monitor',
}

export type ModalDataPayload = {
    [ModalType.ADD_SENSOR]: undefined;
    [ModalType.MODIFY_SENSOR]: Sensor;
    [ModalType.MONITOR_SENSOR]: Sensor;
};

export type OpenedModalData<T extends ModalType> = {
    type: T;
    data: ModalDataPayload[T];
};

export function isModalData<T extends ModalType>(
    type: T,
    modalData: OpenedModalData<ModalType> | null
): modalData is OpenedModalData<T> {
    return modalData?.type === type;
}
