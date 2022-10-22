import type { ID, ParkingLot } from "$lib/types";

export type LayoutData = {
    parkingLots: Record<ID, ParkingLot>;
};