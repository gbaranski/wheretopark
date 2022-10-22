import { env } from '$env/dynamic/public';

const get = (name: string): string => {
    const key = `PUBLIC_${name}`;
    const v = env[key];
    if(v == null && v == undefined) {
        throw new Error(`environment variable ${key} is missing`);
    }
    return v;
}

const MAPBOX_ACCESS_TOKEN = get("MAPBOX_ACCESS_TOKEN");

export {
    MAPBOX_ACCESS_TOKEN,
}