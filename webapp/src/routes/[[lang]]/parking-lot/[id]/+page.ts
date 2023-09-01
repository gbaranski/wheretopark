import { parkingLots } from "$lib/store";
import type { PageLoad } from "./$types";

const waitForParkingLot = (id: string) => new Promise((resolve) => {
    parkingLots.subscribe((parkingLots) => {
        const parkingLot = parkingLots[id];
        if (parkingLot) resolve(parkingLot);
    })
});

export const load = (async ({ params }: { params: { id: string, }}) => {
    const parkingLot = await waitForParkingLot(params.id);
    return {
        parkingLot,
    };
}) satisfies PageLoad;