import { getAllSensors } from "@/helpers/sensor";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { authFetch } from "@/helpers/fetch";

export const load: PageLoad = async ({ fetch }) => {
    const sensors = await getAllSensors(fetch);

    if (sensors.isError) {
        error(500, sensors.error);
    }

    return {
        sensors: sensors.data.map((e) => {
            return { label: e.name, value: e.id };
        }),
    };
};
