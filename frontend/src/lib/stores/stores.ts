import type { ModalType, OpenedModalData } from '@/types/modal';
import type { Sensor } from '@/types/sensor';
import { writable } from 'svelte/store';

export const sensorStore = writable([] as Sensor[]);
export const openedModalStore = writable(
    null as OpenedModalData<ModalType> | null
);
