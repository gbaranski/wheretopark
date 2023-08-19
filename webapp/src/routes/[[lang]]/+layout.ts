import { getParkingLots } from "$lib/client";
import type { LayoutLoad } from "./$types";

export const prerender = false;
export const ssr = false;

export const load = (async ({}) => {
    return {
        parkingLots: await getParkingLots(),
    }
}) satisfies LayoutLoad;