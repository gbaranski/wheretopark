import { distanceBetweenPoints } from "$lib/utils";
import geoip from 'geoip-lite';
import type { LayoutServerLoad } from "./$types";

export const prerender = false;
export const ssr = false;

export const load = (async ({getClientAddress}) => {
    const { getParkingLots } = await import("$lib/server/client");
    const parkingLotsMap = await getParkingLots();
    const userAddress =  getClientAddress();
    const userGeo = geoip.lookup(userAddress);
    if (!userGeo) {
        return {
            parkingLots: parkingLotsMap,
        };
    }
    const userGeoJSON: GeoJSON.Point = {
        type: 'Point',
        coordinates: [userGeo!.ll[0], userGeo!.ll[1]],
    };
    const parkingLots = Object.entries(parkingLotsMap);
    parkingLots.sort((a, b) => 
        distanceBetweenPoints(userGeoJSON, b[1].metadata.geometry) - distanceBetweenPoints(userGeoJSON, a[1].metadata.geometry) 
    );

    return {
        parkingLots: Object.fromEntries(parkingLots),
    }
}) satisfies LayoutServerLoad;