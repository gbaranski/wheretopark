import type {ID, ParkingLot } from "./types";
import { dev } from '$app/environment';
import axios from 'axios';

const serverURL = dev ? "http://localhost:1234" : "https://api.wheretopark.app";
const providersURLs = [
    `${serverURL}/collector`,
    `${serverURL}/cctv`,
];

const getParkingLotFromProvider = async (url: string): Promise<Record<ID, ParkingLot>> => {
    console.log({url});
    const response = await axios.get(`${url}/parking-lots`);
    const parkingLots = response.data;
    console.log({parkingLots});
    return parkingLots;
}

export const getParkingLots = async (): Promise<Record<ID, ParkingLot>> => {
    const requests = await Promise.all(providersURLs.map((url) => getParkingLotFromProvider(url)));
    const parkingLots = Object.assign({}, ...requests);
    console.log({parkingLots});
    return parkingLots;
}