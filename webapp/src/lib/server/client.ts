import {
  APP_ID,
  AWS_API_KEY,
  AWS_REGION,
  CONFIG_PROFILE_ID,
  ENV_ID,
} from "$lib/environment";
import type { ID, ParkingLot } from "../types";
import {
  AppConfigDataClient,
  GetLatestConfigurationCommand,
  type GetLatestConfigurationCommandInput,
  StartConfigurationSessionCommand,
  type StartConfigurationSessionRequest,
} from "@aws-sdk/client-appconfigdata";

type Collector = {
  name: string;
  url: URL;
};

type AppConfiguration = {
  collectors: Collector[];
};

const appConfigClient = new AppConfigDataClient({ region: AWS_REGION });
const getConfigurationToken = async (): Promise<string> => {
  const input: StartConfigurationSessionRequest = {
    ApplicationIdentifier: APP_ID,
    EnvironmentIdentifier: ENV_ID,
    ConfigurationProfileIdentifier: CONFIG_PROFILE_ID,
    RequiredMinimumPollIntervalInSeconds: 60,
  };
  const command = new StartConfigurationSessionCommand(input);
  const response = await appConfigClient.send(command);
  if (response.InitialConfigurationToken == undefined) {
    throw new Error("could not retrieve configuration token");
  }
  return response.InitialConfigurationToken;
};

const getLatestConfiguration = async (): Promise<AppConfiguration> => {
  const configurationToken = await getConfigurationToken();
  const params: GetLatestConfigurationCommandInput = {
    ConfigurationToken: configurationToken,
  };
  const command = new GetLatestConfigurationCommand(params);
  const data = await appConfigClient.send(command);
  if (data.Configuration == undefined) {
    throw new Error("configuration is empty");
  }
  if (data.ContentType != "application/json") {
    throw new Error("configuration is not json");
  }
  const configurationString = new TextDecoder().decode(data.Configuration);
  const configuration = JSON.parse(
    configurationString,
  ) as AppConfiguration;
  return configuration;
};

const getParkingLotsFromCollector = async (
  collector: Collector,
): Promise<ParkingLot[]> => {
  const response = await fetch(`${collector.url}/parking-lots`, {
    headers: {
      "x-api-key": AWS_API_KEY,
    },
  });
  const parkingLots = await response.json();
  return parkingLots;
};

const getParkingLotsS = async (collectors: Collector[]) => {
  const parkingLotRequests = collectors.map(async (collector) => ({
    collector,
    data: await getParkingLotsFromCollector(collector),
  }));
  const parkingLotsByCollector = (await Promise.all(parkingLotRequests)).flat();
  const parkingLots = parkingLotsByCollector.flatMap((e) => e.data);
  return parkingLots;
};

const appConfig = await getLatestConfiguration();
const parkingLots = getParkingLotsS(appConfig.collectors);
// console.log({parkingLots});

export const getParkingLots = async (): Promise<Record<ID, ParkingLot>> => {
  return {};
  // const rawParkingLots = await db.select<ParkingLot & { id?: string}>("parking_lot");
  // const parkingLots = rawParkingLots.map((parkingLot) => {
  //     const id = parkingLot.id!.split(":")[1]
  //     delete parkingLot.id;
  //     return [id, parkingLot];
  // });
  // return Object.fromEntries(parkingLots);
};

export const getParkingLot = async (id: string): Promise<ParkingLot | null> => {
  return null;
  // const rawParkingLot = await db.select<ParkingLot & { id?: string}>(`parking_lot:${id}`);
  // const parkingLot = rawParkingLot[0];
  // delete parkingLot.id;
  // return parkingLot;
};
