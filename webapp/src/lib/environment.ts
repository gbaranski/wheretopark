import { env } from '$env/dynamic/public';

const get = (name: string): string => {
    const key = `PUBLIC_${name}`;
    // @ts-ignore TODO: Remove this ignore
    const v = env[key];
    if(v == null && v == undefined) {
        throw new Error(`environment variable ${key} is missing`);
    }
    return v;
}

const MAPBOX_ACCESS_TOKEN = get("MAPBOX_ACCESS_TOKEN");
const AWS_API_KEY = get("AWS_API_KEY");
const AWS_REGION = "eu-central-1";
const APP_ID = "2cww2oe";
const PROD_ENV_ID = "47osndv";
const ENV_ID = PROD_ENV_ID;
const CONFIG_PROFILE_ID = "8k18g1o";

export {
    MAPBOX_ACCESS_TOKEN,
    AWS_API_KEY,
    AWS_REGION,
    APP_ID,
    ENV_ID,
    CONFIG_PROFILE_ID,
}