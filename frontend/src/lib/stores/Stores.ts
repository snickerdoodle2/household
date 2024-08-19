import type { ModalType, OpenedModalData } from '@/types/Modal.types';
import type { Sensor } from '@/types/Sensor.types';
import { writable } from 'svelte/store';
import { PageType } from '@/types/Page.types';

export const sensorStore = writable([] as Sensor[]);

export const currentPageStore = writable(PageType.SENSOR);

export const openedModalStore = writable(
    null as OpenedModalData<ModalType> | null
);
