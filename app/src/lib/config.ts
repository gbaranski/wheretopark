import { env } from '$env/dynamic/public';

const { PUBLIC_SERVER_URL } = env;
const productionServer = "https://api.wheretopark.app";
export const serverURL = PUBLIC_SERVER_URL ?? productionServer;