/** @type {import('next').NextConfig} */
const nextConfig = {
    output: 'standalone',
    reactStrictMode: true,
    swcMinify: true,
    publicRuntimeConfig: {
        MAPBOX_ACCESS_TOKEN: process.env.MAPBOX_ACCESS_TOKEN
    },
    serverRuntimeConfig: {
        CLIENT_ID: process.env.CLIENT_ID,
        CLIENT_SECRET: process.env.CLIENT_SECRET,
        AUTHORIZATION_URL: process.env.AUTHORIZATION_URL,
        STOREKEEPER_URL: process.env.STOREKEEPER_URL
    }
}

module.exports = nextConfig
