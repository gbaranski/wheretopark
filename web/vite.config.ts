import { sveltekit } from "@sveltejs/kit/vite";
import type { UserConfig } from "vite";
import { SvelteKitPWA } from "@vite-pwa/sveltekit";

const config: UserConfig = {
  plugins: [
    sveltekit(),
    SvelteKitPWA({
      kit: {},
      workbox: {
        globPatterns: ["**/*.{js,css,ico,png,svg,webp,woff,woff2}"],
      },
      registerType: "autoUpdate",
      manifest: {
        short_name: "Where To Park",
        name: "Where To Park: Find A Parking Lot Nearby",
        start_url: "/",
        scope: "/",
        display: "standalone",
        theme_color: "#ffffff",
        background_color: "#ffffff",
        icons: [
          {
            src: "logo_192.png",
            type: "image/png",
            sizes: "192x192",
          },
          {
            src: "logo_512.png",
            type: "image/png",
            sizes: "512x512",
          },
          {
            src: "maskable_logo_192.png",
            type: "image/png",
            sizes: "192x192",
            purpose: "maskable",
          },
          {
            src: "maskable_logo_512.png",
            type: "image/png",
            sizes: "512x512",
            purpose: "maskable",
          },
        ],
      },
    }),
  ],
};

export default config;
