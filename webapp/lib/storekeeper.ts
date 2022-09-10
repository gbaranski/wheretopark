import { Feature, Weekday, type ID, type Metadata, type ParkingLot, type PricingRule, type State } from "./parkingLot";
import axios from 'axios'
import { capitalizeFirstLetter } from "./utils";

const baseURL = "https://storekeeper.wheretopark.app"

const mapMetadata = (v: any): Metadata => ({
    name: v["name"],
    address: v["address"],
    location: {
        latitude: v["location"]["latitude"],
        longitude: v["location"]["longitude"]
    },
    emails: v["emails"],
    phoneNumbers: v["phone-numbers"],
    websites: v["websites"],
    totalSpots: v["total-spots"],
    features: v["features"].map((f: string): Feature => Feature[capitalizeFirstLetter(f) as keyof typeof Feature]),
    currency: v["currency"],
    rules: v["rules"].map((rule: any) => ({
        weekdays: rule["weekdays"] || { start: Weekday.Monday, end: Weekday.Sunday },
        hours: rule["hours"] || null,
        pricing: rule["pricing"].map((pricing: any): PricingRule => ({
            duration: pricing["duration"],
            price: pricing["price"],
            repeating: pricing["repeating"],
        }))
    }))
})

const mapState = (v: any): State => ({
    lastUpdated: v["last-updated"],
    availableSpots: v["available-spots"]
})


export const fetchParkingLotMetadatas = async (): Promise<Record<ID, Metadata>> => {
    const response = await axios({
        method: 'get',
        url: `${baseURL}/parking-lot/metadata`,
        responseType: 'json',
    });
    const metadatas = Object.entries(response.data).map(([id, v]): [ID, Metadata] => ([id, mapMetadata(v)]))
    return Object.fromEntries(metadatas)
}


export const fetchParkingLotStates = async (): Promise<Record<ID, State>> => {
    const response = await axios({
        method: 'get',
        url: `${baseURL}/parking-lot/state`,
        responseType: 'json',
    });
    const metadatas = Object.entries(response.data).map(([id, v]): [ID, State] => ([id, mapState(v)]))
    return Object.fromEntries(metadatas)
}


export const fetchParkingLots = async (): Promise<ParkingLot[]> => {
    const [metadatas, states] = await Promise.all([fetchParkingLotMetadatas(), fetchParkingLotStates()]);
    const parkingLots = Object.entries(metadatas).map(([id, metadata]): ParkingLot => {
        return {
            id: id,
            metadata,
            state: states[id]
        }
    })
    return parkingLots
}


export const fetchParkingLotMetadata = async (id: ID): Promise<Metadata> => {
    const response = await axios({
        method: 'get',
        url: `${baseURL}/parking-lot/${id}/metadata`,
        responseType: 'json',
    });
    const metadata = mapMetadata(response.data);
    return metadata;
}


export const fetchParkingLotState = async (id: ID): Promise<State> => {
    const response = await axios({
        method: 'get',
        url: `${baseURL}/parking-lot/${id}/state`,
        responseType: 'json',
    });
    const state = mapState(response.data);
    return state;
}


export const fetchParkingLot = async (id: ID): Promise<ParkingLot> => {
    const [metadata, state] = await Promise.all([fetchParkingLotMetadata(id), fetchParkingLotState(id)]);
    return {
        id,
        metadata,
        state
    }
}