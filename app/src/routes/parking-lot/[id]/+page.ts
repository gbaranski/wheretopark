import type { ParkingLot } from "$lib/parkingLot";
import { parkingLots } from "$lib/store";
import type { PageLoad } from "./$types";
import { building } from '$app/environment';
import { getForecast } from "$lib/forecaster";
import dayjs from "dayjs";


const isCapacitor = building && process.env.CAPACITOR == "true";
export const ssr = isCapacitor ? false : true;
export const prerender = false;

const waitForParkingLot = (id: string): Promise<ParkingLot> => new Promise((resolve) => {
    parkingLots.subscribe((parkingLots) => {
        const parkingLot = parkingLots[id];
        if (parkingLot) resolve(parkingLot);
    })
});

export const load = (async ({ fetch, params }) => {
    const parkingLot = await waitForParkingLot(params.id);
    const forecast = await getForecast(params.id, dayjs().format("YYYY-MM-DD"), fetch)
    return {
        parkingLot,
        forecast,
    };
}) satisfies PageLoad;