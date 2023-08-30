import { updateParkingLots } from "$lib/client";
import type { LayoutLoad } from "./$types";

export const prerender = true;
export const ssr = false;

export const load = (async ({fetch}) => {
    updateParkingLots(fetch);
}) satisfies LayoutLoad;