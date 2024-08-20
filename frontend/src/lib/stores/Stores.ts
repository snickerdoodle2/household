import type { ModalType, OpenedModalData } from '@/types/Modal.types';
import type { Sensor } from '@/types/Sensor.types';
import { writable } from 'svelte/store';
import { PageType } from '@/types/Page.types';
import { mapStore } from './MapStore';

export const sensorStore = writable([] as Sensor[]);

export const sensorValueMap = writable(
    [] as { id: Sensor['id']; val: number }[]
);

export const currentPageStore = writable(PageType.SENSOR);

export const openedModalStore = writable(
    null as OpenedModalData<ModalType> | null
);
