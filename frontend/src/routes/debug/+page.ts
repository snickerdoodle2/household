import { getAllSensors } from "@/helpers/sensor";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { authFetch } from "@/helpers/fetch";

export const load: PageLoad = async () => {
    const sensors = await getAllSensors(authFetch);

    if (sensors.isError) {
        error(500, sensors.error);
    }

    return {
        sensors: sensors.data.map((e) => {
            return { label: e.name, value: e.id };
        }),
    };
};
