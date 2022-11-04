import { getParkingLot } from "$lib/server/client";
import type { LayoutServerLoad } from ".svelte-kit/types/src/routes/$types";

export const load: LayoutServerLoad = async ({ params }: { params: { id?: string, }}) => {
    if (params.id == null) throw new Error("Missing id in params");
    const parkingLot = await getParkingLot(params.id)!;
    return {
        parkingLot: parkingLot
    };
}