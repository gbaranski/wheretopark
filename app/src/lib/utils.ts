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

export const googleMapsLink = (geometry: GeoJSON.Point) => {
  const [longitude, latitude] = geometry.coordinates;
  return `https://google.com/maps/place/${latitude},${longitude}`;
};

export const distanceBetweenPoints = (a: GeoJSON.Point, b: GeoJSON.Point) => {
  return haversine({ geometry: a }, { geometry: b }, { format: "geojson" });
};
