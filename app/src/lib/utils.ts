import haversine from "haversine";

export function capitalizeFirstLetter(s: string) {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

// import {convertDistance} from 'geolib';
// import {Coordinate} from './types';

// export const prettyDistance = (from: Coordinate, to: Coordinate) => {
//     console.log({from, to});
//     const distance = from.distanceTo(to);
//     if (distance < 1000) {
//         return `${Math.round(distance / 100) * 100} m`;
//     } else {
//         return `${Math.round(convertDistance(distance, 'km') * 10) / 10} km`;
//     }
// }

import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import duration from "dayjs/plugin/duration";
import utc from "dayjs/plugin/utc";
import timezone from "dayjs/plugin/timezone";
import localeData from "dayjs/plugin/localeData";
import { Capacitor } from "@capacitor/core";
import {
  Geolocation as NativeGeolocation,
  type Position as NativePosition,
} from "@capacitor/geolocation";
import type { Thing, WithContext } from "schema-dts";

dayjs.extend(localeData);
dayjs.extend(relativeTime);
dayjs.extend(duration);
dayjs.extend(utc);
dayjs.extend(timezone);

export const humanizeDuration = (s: string): string => {
  return dayjs.duration(s).humanize(false);
};

export const weekdays = dayjs.weekdays();
weekdays.entries;

export const getWeekday = (): number => {
  return dayjs().day();
};

export const googleMapsLink = (geometry: GeoJSON.Point): URL => {
  const [longitude, latitude] = geometry.coordinates;
  return new URL(`https://google.com/maps/place/${latitude},${longitude}`);
};

export const distanceBetweenPoints = (a: GeoJSON.Point, b: GeoJSON.Point) => {
  return haversine({ geometry: a }, { geometry: b }, { format: "geojson" });
};

const positionFromNativeGeolocation = (
  position: NativePosition,
): GeolocationPosition => {
  return {
    ...position,
    coords: {
      ...position.coords,
      altitudeAccuracy: position.coords.altitudeAccuracy || null,
    },
  };
};

export const geolocation = () =>
  Capacitor.isNativePlatform()
    ? nativeGeolocation
    : window.navigator.geolocation;

export const nativeGeolocation: Geolocation = {
  clearWatch: function (watchId: number): void {
    NativeGeolocation.clearWatch({
      id: watchId.toString(),
    });
  },
  getCurrentPosition: async function (
    successCallback: PositionCallback,
    errorCallback?: PositionErrorCallback | null | undefined,
    options?: PositionOptions | undefined,
  ): Promise<void> {
    try {
      const permissions = await NativeGeolocation.checkPermissions();
      console.log({ permissions });
      const currentPosition = await NativeGeolocation.getCurrentPosition(
        options,
      );
      successCallback(positionFromNativeGeolocation(currentPosition));
    } catch (e) {
      console.log("error when getting current position", e);
      if (errorCallback) {
        errorCallback({
          code: 1,
          message: e as string,
          PERMISSION_DENIED: 1,
          POSITION_UNAVAILABLE: 2,
          TIMEOUT: 3,
        });
      }
    }
  },
  watchPosition: function (
    successCallback: PositionCallback,
    errorCallback?: PositionErrorCallback | null,
    options?: PositionOptions,
  ): number {
    const wait = NativeGeolocation.watchPosition(options || {}, (position) => {
      if (position) successCallback(positionFromNativeGeolocation(position));
      else console.error("received null position from native");
    });
    wait.then((watchId) => console.log({ watchId }));
    return 1;
  },
};


export const serializeSchema = <T extends Thing>(thing: WithContext<T>): string => {
  return `<script type="application/ld+json">${JSON.stringify(
    thing,
    null,
    2
  )}</script>`
}