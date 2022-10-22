import type { PageParentData } from "./$types";

export const load = async ({ params, parent }: { params: { id: string, }, parent: () => Promise<PageParentData> }) => {
    const { parkingLots }: PageParentData = await parent();
    const parkingLot = parkingLots[params.id];
    return {
        parkingLot: parkingLot
    };
}