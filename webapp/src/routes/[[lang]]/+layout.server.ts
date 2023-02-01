import type { LayoutServerLoad } from "./$types";
import type { LayoutData } from "$types/layout";
import { getParkingLots } from "$lib/server/client";
import { distanceBetweenPoints } from "$lib/utils";
import geoip from 'geoip-lite';

export const load: LayoutServerLoad = async ({getClientAddress}): Promise<LayoutData> => {
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
}