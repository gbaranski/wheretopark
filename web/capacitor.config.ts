import { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
  appId: 'com.gbaranski.wheretopark',
  appName: 'Where To Park',
  webDir: 'build',
  server: {
    androidScheme: 'https'
  }
};

export default config;
