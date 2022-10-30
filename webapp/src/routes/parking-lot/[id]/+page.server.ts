import { getParkingLot } from "$lib/server/client";

export const load = async ({ params }: { params: { id: string, }}) => {
    const parkingLot = await getParkingLot(params.id)!;
    return {
        parkingLot: parkingLot
    };
}