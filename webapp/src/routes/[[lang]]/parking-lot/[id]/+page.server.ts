import type { PageServerLoad } from "./$types";

export const load = (async ({ params }: { params: { id?: string, }}) => {
    const { getParkingLots } = await import("$lib/client");
    const parkingLots = await getParkingLots();

    return {
        parkingLot: parkingLots[params.id!]
    };
}) satisfies PageServerLoad;