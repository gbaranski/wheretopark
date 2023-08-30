import { writable } from 'svelte/store';
import type { ID, ParkingLot } from './types';

export const currentMap = writable<mapboxgl.Map | null>(null);
export const parkingLots = writable<Record<ID, ParkingLot>>({});
