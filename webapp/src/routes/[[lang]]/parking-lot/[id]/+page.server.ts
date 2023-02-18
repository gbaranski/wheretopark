import type { PageServerLoad } from "./$types";

export const load = (async ({ params }: { params: { id?: string, }}) => {
    const { getParkingLot } = await import("$lib/server/client");

    if (params.id == null) throw new Error("Missing id in params");
    const parkingLot = await getParkingLot(params.id)!;
    return {
        parkingLot: parkingLot
    };
}) satisfies PageServerLoad;