import { defineConfig } from 'astro/config';

import tailwind from "@astrojs/tailwind";

// https://astro.build/config
export default defineConfig({
  integrations: [tailwind()],
  experimental: {
    i18n: {
      defaultLocale: "en",
      locales: ["en", "pl"],
      routingStrategy: "prefix-other-locales",
    }
  },
});