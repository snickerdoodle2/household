import { SERVER_URL } from "$lib/const";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { z } from "zod";
import { Sensor } from "@/types/sensor";

export const load: PageLoad = async ({ fetch }) => {
    const res = await fetch(`${SERVER_URL}/api/v1/sensor`);
    const data = await res.json();
    if (!res.ok) {
        error(res.status, data);
    }

    const parsed = z.object({ data: Sensor.array() }).safeParse(data);
    if (!parsed.success) {
        console.log(parsed.error);
        error(500, { message: "Error parsing the data" });
    }

    return {
        sensors: parsed.data.data.map((e) => {
            return { label: e.name, value: e.id };
        }),
    };
};
