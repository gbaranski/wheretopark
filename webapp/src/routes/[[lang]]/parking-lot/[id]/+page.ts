import type { PageLoad } from "./$types";
import { getParkingLots } from "$lib/client";

export const load = (async ({ params }: { params: { id?: string, }}) => {
    const parkingLots = await getParkingLots();

    return {
        parkingLot: parkingLots[params.id!]
    };
}) satisfies PageLoad;