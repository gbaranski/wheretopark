import { writable } from 'svelte/store';
import type { Feature, ID, ParkingLot } from './parkingLot';

export const currentMap = writable<mapboxgl.Map | null>(null);
export const parkingLots = writable<Record<ID, ParkingLot>>({});
export const isLoading = writable<boolean>(true);


export type SearchFilters = {
    openNow: boolean;
    hasFeatures: Record<Feature, boolean>
    minAvailableSpots: number;
}

export const searchFilters = writable<SearchFilters>({
    openNow: false,
    hasFeatures: {
        UNCOVERED: false,
        COVERED: false,
        UNDERGROUND: false,
        GUARDED: false,
    },
    minAvailableSpots: 0
});