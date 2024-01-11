import { serverURL } from "./config";
import type { ID } from "./parkingLot";

export type Prediction = {
    date: Date;
    occupancy: number;
}

export type Forecast = {
    predictions: Prediction[];
}

export const getForecast = async (
    id: ID,
    date: string,
    fetch: typeof window.fetch,
): Promise<Forecast> => {
    const url = `${serverURL}/v1/forecaster/predict/${id}/${date}`;
    const response = await fetch(url);
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();
    return {
        predictions: data["predictions"].map((p: any) => ({
            date: new Date(p.date),
            occupancy: p.occupancy,
        })),
    };
}
