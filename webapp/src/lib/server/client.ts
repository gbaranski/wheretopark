import type {ID, Metadata, ParkingLot, State } from "../types";
import {DATABASE_PASSWORD, DATABASE_URL, DATABASE_USER} from "./secrets";
import Surreal from 'surrealdb.js'

const db = new Surreal(DATABASE_URL);
await db.signin({
    user: DATABASE_USER,
    pass: DATABASE_PASSWORD,
});
const isDevelopment = process.env.NODE_ENV === 'development';
await db.use("wheretopark", isDevelopment ? "development" : "production");

export const getParkingLots = async (): Promise<Record<ID, ParkingLot>> => {
    const metadatas = await db.select<Metadata & { id: string}>("metadata");
    const states = await db.select<State & {id: string}>("state");
    const parkingLots = metadatas.map((metadata) => {
        const id = metadata.id.split(":")[1]
        const state = states.find((state) => state.id == `state:${id}`);
        return [id, {
            metadata,
            state,
        } as ParkingLot];
    });
    return Object.fromEntries(parkingLots);
}