import type {ID, ParkingLot } from "./types";
import { parkingLots as parkingLotsStore } from "$lib/store";
import { dev } from '$app/environment';

const serverURL = dev ? "http://localhost:1234" : "https://api.wheretopark.app";
const providersURLs = [
    `${serverURL}/collector`,
    // `${serverURL}/cctv`,
];

const getParkingLotFromProvider = async function* (fetch: typeof window.fetch, url: string): AsyncGenerator<Record<ID, ParkingLot>> {
    const response = await fetch(`${url}/parking-lots`);
    const reader = response.body!.getReader();
    let buffer = '';
    const decoder = new TextDecoder("utf-8")
    while (true) {
        const { done, value } = await reader.read();
        if (done) {
            if (buffer) {
                const jsonObject = JSON.parse(buffer);
                yield jsonObject;
            }
            break;
        }
        buffer += decoder.decode(value, {stream: true});
        // Split by newline to process complete JSON objects
        const parts = buffer.split('\r\n');
        buffer = parts.pop()!;

        for (const part of parts) {
            if (part.trim() !== '') {
                const jsonObject = JSON.parse(part);
                yield jsonObject;
            }
        }
    }
}

export const updateParkingLots = async (fetch: typeof window.fetch) => {
    const promises = providersURLs.map(async (url) => {
        const providerParkingLots = getParkingLotFromProvider(fetch, url);
        for await (const parkingLots of providerParkingLots) {
            parkingLotsStore.update((value) => {
                return {...value, ...parkingLots};
            })
        }
    });
    await Promise.all(promises);
}