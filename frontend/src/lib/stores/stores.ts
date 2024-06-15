import type { Sensor } from "@/types/sensor";
import { writable } from "svelte/store";

export const sensors = writable([] as Sensor[]);