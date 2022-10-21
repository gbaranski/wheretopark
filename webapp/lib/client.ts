import getConfig from "next/config";
import {AuthorizationClient, StorekeeperClient} from "./types";
import {app} from "wheretopark-shared";
import AccessType = app.wheretopark.shared.AccessType;

const {serverRuntimeConfig} = getConfig();
console.log({serverRuntimeConfig});
const authorizationClient = new AuthorizationClient(serverRuntimeConfig.AUTHORIZATION_URL, serverRuntimeConfig.CLIENT_ID, serverRuntimeConfig.CLIENT_SECRET!)
export const storekeeperClient = new StorekeeperClient(
    serverRuntimeConfig.STOREKEEPER_URL,
    [AccessType.ReadMetadata, AccessType.ReadState],
    authorizationClient
)
