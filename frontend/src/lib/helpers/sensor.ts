import { SERVER_URL } from "@/const";
import type { Result } from "@/types/result";
import { SensorSchema, type SensorType } from "@/types/sensor";
import { z } from "zod";

export const getAllSensors = async (
    fetch: (
        input: RequestInfo | URL,
        init?: RequestInit | undefined,
    ) => Promise<Response>,
): Promise<Result<SensorType[], string>> => {
    const res = await fetch(`${SERVER_URL}/api/v1/sensor`);
    const data = await res.json();
    if (!res.ok) {
        return {
            isError: true,
            error: data.error,
        };
    }

    const parsed = z.object({ data: SensorSchema.array() }).safeParse(data);
    if (!parsed.success) {
        return {
            isError: true,
            error: "Error while parsing the data! (getAllSensors)",
        };
    }

    return {
        isError: false,
        data: parsed.data.data,
    };
};
