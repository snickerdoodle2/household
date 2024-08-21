import { sensorStore, sensorValueMap as sensorValues } from '@/stores/Stores';
import { getAllSensors, getSensorData } from './requests/Sensor.requests';
import { get } from 'svelte/store';
import type { Sensor } from '@/types/Sensor.types';
import { syncCategories } from './requests/Categories.requests';

export async function initializeStores() {
    await syncSensorConfig();
    await syncSensorValues();
    await syncCategories();
    console.log('Sensor data initialized!');
}

export async function syncSensorConfig() {
    const response = await getAllSensors(fetch);
    if (response.isError) {
        console.error('Failed to fetch sensors!', response.error);
        return;
    }
    sensorStore.set(response.data);
    console.info('Sensors updated:', response.data);
}

export async function syncSensorValues() {
    const values: { id: Sensor['id']; val: number }[] = [];
    for (const sensor of get(sensorStore)) {
        const value = await getSensorData(sensor.id);
        values.push({ id: sensor.id, val: value });
    }
    sensorValues.set(values);
    console.info('Sensor values updated: ', get(sensorValues));
}
