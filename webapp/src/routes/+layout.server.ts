import type { LayoutServerLoad } from "./$types";
import type { LayoutData } from "src/types/layout";
import { getParkingLots } from "$lib/server/client";

export const load: LayoutServerLoad = async (): Promise<LayoutData> => {
    const parkingLots = await getParkingLots();
    return {
        parkingLots: parkingLots,
    }
}