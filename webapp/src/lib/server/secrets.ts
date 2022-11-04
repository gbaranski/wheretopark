import { env } from '$env/dynamic/private';

const get = (name: string): string => {
    const key = `SECRET_${name}`;
    const v = env[key];
    if(v == null && v == undefined) {
        throw new Error(`environment variable ${key} is missing`);
    }
    return v;
}

const DATABASE_USER = get("DATABASE_USER");
const DATABASE_PASSWORD = get("DATABASE_PASSWORD");
const DATABASE_URL = get("DATABASE_URL");
const DATABASE_NAME = get("DATABASE_NAME");

export {
    DATABASE_USER,
    DATABASE_PASSWORD,
    DATABASE_URL,
    DATABASE_NAME,
}