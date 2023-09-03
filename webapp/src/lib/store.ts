import { writable } from 'svelte/store';
import type { ID, ParkingLot } from './parkingLot';

export const currentMap = writable<mapboxgl.Map | null>(null);
export const parkingLots = writable<Record<ID, ParkingLot>>({});
export const isLoading = writable<boolean>(true);