import { writable } from 'svelte/store';

export const currentMap = writable<mapboxgl.Map | null>(null);