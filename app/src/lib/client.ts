import { type ID, ParkingLot } from "./parkingLot";
import { isLoading, parkingLots as parkingLotsStore } from "$lib/store";
import { dev } from "$app/environment";

const serverURL = dev ? "http://localhost:1234" : "https://api.wheretopark.app";

type Provider = {
  name: string;
  url: URL;
};

const getProviders = async (
  fetch: typeof window.fetch,
): Promise<Provider[]> => {
  const response = await fetch(`${serverURL}/v1/providers`);
  const providers = await response.json() as Provider[];
  return providers;
};

const parse = (s: string): Record<ID, ParkingLot> => {
  const rawParkingLots = JSON.parse(s) as Record<ID, any>;
  const parkingLots = Object.entries(rawParkingLots).map(([id, raw]) => {
    const parkingLot = ParkingLot.fromJSON(id, raw);
    return [id, parkingLot];
  });
  return Object.fromEntries(parkingLots);
};

const getParkingLotFromProvider = async function* (
  fetch: typeof window.fetch,
  provider: Provider,
): AsyncGenerator<Record<ID, ParkingLot>> {
  const response = await fetch(`${provider.url}/parking-lots`);
  const reader = response.body!.getReader();
  let buffer = "";
  const decoder = new TextDecoder("utf-8");
  while (true) {
    const { done, value } = await reader.read();
    if (done) {
      if (buffer) yield parse(buffer);
      break;
    }
    buffer += decoder.decode(value, { stream: true });
    // Split by newline to process complete JSON objects
    const parts = buffer.split("\r\n");
    buffer = parts.pop()!;

    for (const part of parts) {
      if (part.trim() !== "") {
        yield parse(part);
      }
    }
  }
};

export const getParkingLots = async (fetch: typeof window.fetch): Promise<Record<ID, ParkingLot>> => {
  const allParkingLots: Record<ID, ParkingLot> = {};
  const providers = await getProviders(fetch);
  const promises = providers.map(async (provider) => {
    const providerParkingLots = getParkingLotFromProvider(fetch, provider);
    for await (const parkingLots of providerParkingLots) {
      Object.entries(parkingLots).forEach(([id, parkingLot]) => {
        allParkingLots[id] = parkingLot;
      });
    }
  });
  await Promise.all(promises);
  return allParkingLots;
}

export const updateParkingLots = async (fetch: typeof window.fetch) => {
  const providers = await getProviders(fetch);
  const promises = providers.map(async (provider) => {
    const providerParkingLots = getParkingLotFromProvider(fetch, provider);
    for await (const parkingLots of providerParkingLots) {
      parkingLotsStore.update((value) => {
        return { ...value, ...parkingLots };
      });
    }
  });
  await Promise.all(promises);
  isLoading.set(false);
};
