import { error, type RequestHandler } from "@sveltejs/kit";
import { spawnSync } from "child_process";

export const GET: RequestHandler = async ({ url }) => {
    const src = url.searchParams.get("src");
    if (src == null) {
        throw error(400, "missing src")
    }
    const srcUrl = new URL(src);
    const image = await fetchImage(srcUrl);
    // fs.writeFileSync("output.jpg", image);
    return new Response(image, { headers: { "Content-Type": "image/jpeg" } });
}

const fetchImage = async (src: URL): Promise<Buffer> => {
    const args = ["-i", src.toString(), "-frames:v", "1", "-vf", "select=gte(n\\,5)", "-f", "image2pipe", "-"];
    if (src.protocol === "rtsp:") {
        args.unshift("-rtsp_transport", "tcp", "-buffer_size", "2048");
    }
    const child = spawnSync("ffmpeg", args);
    if (child.status !== 0) {
        throw new Error(`ffmpeg exited with code ${child.status}. output: ${child.stderr}`);
    }
    return child.stdout;
}