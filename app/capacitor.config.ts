import { CapacitorConfig } from "@capacitor/cli";

const liveReload = process.env.LIVE_RELOAD == "true";

const config: CapacitorConfig = {
  appId: "com.gbaranski.wheretopark",
  appName: "Where To Park",
  webDir: "build",
  server: {
    androidScheme: "https",
    url: liveReload ? "http://192.168.1.2:5173" : undefined,
    cleartext: liveReload ? true : undefined,
  },
  backgroundColor: "#ffffff",
  ios: {
    allowsLinkPreview: false,
  }
};

export default config;
