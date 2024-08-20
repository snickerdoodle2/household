import { writable, type Writable } from 'svelte/store';

export function mapStore<K, V>(initial: [K, V][] = []) {
    const map: Map<K, V> = new Map(initial);
    const { subscribe, set, update }: Writable<Map<K, V>> = writable(map);

    return {
        subscribe,
        set: (key: K, value: V) =>
            update((map) => {
                map.set(key, value);
                return map;
            }),
        delete: (key: K) =>
            update((map) => {
                map.delete(key);
                return map;
            }),
        clear: () =>
            update((map) => {
                map.clear();
                return map;
            }),
        get: (key: K): V | undefined => {
            let result: V | undefined;
            subscribe((map) => (result = map.get(key)))();
            return result;
        },
        size: (): number => {
            let size = 0;
            subscribe((map) => (size = map.size))();
            return size;
        },
        entries: (): [K, V][] => {
            let entries: [K, V][] = [];
            subscribe((map) => (entries = Array.from(map.entries())))();
            return entries;
        },
        setAll: (newMap: [K, V][]) => set(new Map(newMap)),
    };
}
