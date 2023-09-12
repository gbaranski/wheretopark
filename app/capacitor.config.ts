import { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
  appId: 'com.gbaranski.wheretopark',
  appName: 'Where To Park',
  webDir: 'build',
  server: {
    androidScheme: 'https',
    url: "http://192.168.1.2:5173",
    cleartext: true,
  },
  backgroundColor: "#ffffff",
  ios: {
    backgroundColor: "#ffffff",
    contentInset: "always",
    scrollEnabled: true,
  }
};

export default config;
