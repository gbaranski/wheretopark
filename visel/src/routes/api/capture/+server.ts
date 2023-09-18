import { error, type RequestHandler } from "@sveltejs/kit";
import { spawnSync } from "child_process";

export const GET: RequestHandler = async ({ url }) => {
    const src = url.searchParams.get("src");
    if (src == null) {
        throw error(400, "missing src")
    }
    const image = await fetchAndParseHLS(src);
    // fs.writeFileSync("output.jpg", image);
    return new Response(image, { headers: { "Content-Type": "image/jpeg" } });
}

const fetchAndParseHLS = async (src: string): Promise<Buffer> => {
    const child = spawnSync("ffmpeg", ["-i", src, "-frames:v", "1", "-f", "image2pipe", "-"]);
    if (child.status !== 0) {
        throw new Error(`ffmpeg exited with code ${child.status}. output: ${child.stderr}`);
    }
    return child.stdout;
}