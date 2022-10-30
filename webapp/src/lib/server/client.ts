import type {ID, Metadata, ParkingLot, State } from "../types";
import {DATABASE_PASSWORD, DATABASE_URL, DATABASE_USER} from "./secrets";
import Surreal from 'surrealdb.js'

const db = new Surreal(`${DATABASE_URL}/rpc`);
await db.signin({
    user: DATABASE_USER,
    pass: DATABASE_PASSWORD,
});
const isDevelopment = process.env.NODE_ENV === 'development';
await db.use("wheretopark", isDevelopment ? "development" : "production");

export const getParkingLots = async (): Promise<Record<ID, ParkingLot>> => {
    const rawParkingLots = await db.select<ParkingLot & { id?: string}>("parking_lot");
    const parkingLots = rawParkingLots.map((parkingLot) => {
        const id = parkingLot.id!.split(":")[1]
        delete parkingLot.id;
        return [id, parkingLot];
    });
    return Object.fromEntries(parkingLots);
}

export const getParkingLot = async (id: string): Promise<ParkingLot | null> => {
    const rawParkingLot = await db.select<ParkingLot & { id?: string}>(`parking_lot:${id}`);
    const parkingLot = rawParkingLot[0];
    delete parkingLot.id;
    return parkingLot;
}