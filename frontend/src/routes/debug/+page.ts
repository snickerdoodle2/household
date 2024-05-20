import { SERVER_URL } from "$lib/const";
import type { Sensor } from "@/types/sensor";
import type { PageLoad } from "./$types";

export const load: PageLoad = async () => {
    return {
        ids: (await fetch(`${SERVER_URL}/api/v1/sensor`).then(e => e.json())) as { data: Sensor[] }
    }
}
