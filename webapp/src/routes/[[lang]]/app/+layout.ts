import { updateParkingLots } from "$lib/client";
import type { LayoutLoad } from "./$types";

export const load = (async ({fetch}) => {
    updateParkingLots(fetch);
}) satisfies LayoutLoad;