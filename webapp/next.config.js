const getEnv = (key) => {
    const value = process.env[key];
    if (value == undefined)
        throw new Error(`Environment variable ${key} is not set`);
    return value;
}

/** @type {import('next').NextConfig} */
const nextConfig = {
    output: 'standalone',
    reactStrictMode: true,
    swcMinify: true,
    publicRuntimeConfig: {
        MAPBOX_ACCESS_TOKEN: getEnv("MAPBOX_ACCESS_TOKEN")
    },
    serverRuntimeConfig: {
        CLIENT_ID: getEnv("CLIENT_ID"),
        CLIENT_SECRET: getEnv("CLIENT_SECRET"),
        AUTHORIZATION_URL: getEnv("AUTHORIZATION_URL"),
        STOREKEEPER_URL: getEnv("STOREKEEPER_URL")
    }
}

module.exports = nextConfig
