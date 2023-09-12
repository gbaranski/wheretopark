import { getParkingLots } from "$lib/client";
import type { RequestHandler } from "@sveltejs/kit";

export const GET = (async ({ fetch }) => {
  const parkingLots = await getParkingLots(fetch);
  const sitemap = `
<?xml version="1.0" encoding="UTF-8" ?>
<urlset
    xmlns="https://www.sitemaps.org/schemas/sitemap/0.9"
    xmlns:xhtml="https://www.w3.org/1999/xhtml"
    xmlns:mobile="https://www.google.com/schemas/sitemap-mobile/1.0"
    xmlns:news="https://www.google.com/schemas/sitemap-news/0.9"
    xmlns:image="https://www.google.com/schemas/sitemap-image/1.1"
    xmlns:video="https://www.google.com/schemas/sitemap-video/1.1"
>
    ${
    Object.values(parkingLots).map((parkingLot) => {
      const link = parkingLot.link();
      const lastUpdated = parkingLot.lastUpdated.format("YYYY-MM-DD");
      return `
    <url>
        <loc>${link}</loc>
        <lastmod>${lastUpdated}</lastmod>
    </url>`;
    })
  }
</urlset>`;
  return new Response(sitemap.trim(), {
    headers: {
      "Content-Type": "application/xml",
    },
  });
}) satisfies RequestHandler;
