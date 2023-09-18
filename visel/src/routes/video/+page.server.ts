export function load({ url }: {url: URL}) {
    const strURL = url.searchParams.get("src");
    if (!strURL) throw new Error("no src parameter provided")
    if (!URL.canParse(strURL)) throw new Error("invalid src parameter provided")

    return {
        src: strURL
    }
}