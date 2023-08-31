import {
  type Coordinate,
  Feature,
  type LanguageCode,
  type Metadata,
  type ParkingLot,
  ParkingLotStatus,
  type Rule,
  SpotType,
} from "./types";
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

const OF_CATEGORY_FEATURES = [
  Feature.UNCOVERED,
  Feature.COVERED,
  Feature.UNDERGROUND,
];

export const getCategory = (features: Feature[]): string => {
  const ofCategoryFeatures = features.filter((f) =>
    OF_CATEGORY_FEATURES.includes(f)
  );
  if (ofCategoryFeatures.length == 0) {
    return "Unknown";
  } else if (ofCategoryFeatures.length == 1) {
    const category = ofCategoryFeatures[0]!;
    if (category == Feature.COVERED) return "Covered";
    else if (category == Feature.UNDERGROUND) return "Underground";
    else if (category == Feature.UNCOVERED) return "Not covered";
  } else {
    if (
      features.includes(Feature.COVERED) ||
      features.includes(Feature.UNDERGROUND)
    ) {
      if (features.includes(Feature.UNCOVERED)) return "Partially covered";
      else return "Covered & Underground";
    }
  }
  return "Unknown";
};

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

export const timeFromNow = (iso: string) => {
  return dayjs(iso).fromNow();
};

export const humanizeDuration = (s: string): string => {
  return dayjs.duration(s).humanize(false);
};

export const weekdays = dayjs.weekdays();
weekdays.entries;

export const getWeekday = (): number => {
  return dayjs().day();
};

import OpeningHours, { type argument_hash } from "opening_hours";

export const rulesForDay = (
  rules: Rule[],
  spotType: SpotType,
  day: number,
): (Rule & {humanHours: string})[] => {
  const weekday = dayjs().day(day);
  const date = weekday.toDate();
  return rules
    .filter((rule) =>
      rule.applies === undefined || rule.applies.includes(spotType)
    )
    .map((rule) => ({ rule, openingHours: new OpeningHours(rule.hours) }))
    .filter(({ rule, openingHours }) => openingHours.getState(date))
    .map(({ rule, openingHours }) => {
      const matchingRuleIndex = openingHours.getMatchingRule(date);
      console.assert(matchingRuleIndex != undefined);
      //@ts-ignore
      const humanHours = openingHours.prettifyValue({ rule_index: matchingRuleIndex! });
      return { humanHours, ...rule }
    })
};

export const parkingLotStatus = (
  parkingLot: ParkingLot,
  spotType: SpotType,
): [ParkingLotStatus, string?] => {
  const rules = parkingLot.metadata.rules.filter((rule) =>
    rule.applies === undefined || rule.applies.includes(spotType)
  );
  const allDayRule = rules.find((rule) => rule.hours == "24/7");
  if (allDayRule != null) {
    return [ParkingLotStatus.Open, "24/7"];
  }
  const rawOpeningHours = rules.map((rule) => rule.hours).join(", ");
  const openingHours = new OpeningHours(rawOpeningHours);
  const currentDate = dayjs().tz(parkingLot.metadata.timezone);

  const rawNextChange = openingHours.getNextChange(currentDate.toDate());
  const nextChange = rawNextChange != undefined
    ? dayjs(rawNextChange)
    : undefined;
  const hoursToChange = nextChange
    ? nextChange.diff(currentDate, "hours")
    : undefined;
  if (openingHours.getState(currentDate.toDate())) {
    if (hoursToChange == undefined) return [ParkingLotStatus.Open, "24/7"];
    else if (hoursToChange > 1) {
      return [ParkingLotStatus.Open, `Closes ${nextChange!.format("HH:mm")}`];
    } else {
      const matchingRuleIndex = openingHours.getMatchingRule(
        currentDate.toDate(),
      );
      const argumentHash = {
        // @ts-ignore
        rule_index: matchingRuleIndex! as "number",
      } as Partial<argument_hash>;
      // @ts-ignore
      const matchingRule = openingHours.prettifyValue(argumentHash);
      return [ParkingLotStatus.ClosesSoon, matchingRule];
    }
  } else {
    if (hoursToChange == undefined) {
      return [ParkingLotStatus.Closed, "Temporarily closed"];
    }
    if (hoursToChange < 1) {
      return [
        ParkingLotStatus.OpensSoon,
        `Open ${nextChange!.format("HH:mm")}`,
      ];
    } else {return [
        ParkingLotStatus.Closed,
        `Opens ${nextChange!.format("dd HH:mm")}`,
      ];}
  }
};

export const parkingLotStatusColor = (status: ParkingLotStatus) => {
  switch (status) {
    case ParkingLotStatus.Open:
      return "text-success";
    case ParkingLotStatus.ClosesSoon:
      return "text-warning";
    case ParkingLotStatus.OpensSoon:
      return "text-warning";
    case ParkingLotStatus.Closed:
      return "text-error";
  }
};

export const resourceText = (resource: URL) => {
  switch (resource.protocol) {
    case "http:":
    case "https:":
      return resource.host;
    case "mailto:":
      return resource.pathname;
    case "tel:":
      return resource.pathname.replaceAll("-", " ");
    default:
      return "";
  }
};

export const spotTypeIcon = (spotType: SpotType) => {
  console.log({ spotType });
  switch (spotType) {
    case SpotType.CAR:
      return "directions_car";
    case SpotType.CAR_DISABLED:
      return "accessible";
    case SpotType.CAR_ELECTRIC:
      return "electric_bolt";
    case SpotType.MOTORCYCLE:
      return "motorcycle";
    case SpotType.TRUCK:
      return "local_shipping";
    case SpotType.BUS:
      return "directions_bus";
    default:
      return "error_outline";
  }
};

export const googleMapsLink = (geometry: GeoJSON.Point) => {
  const [longitude, latitude] = geometry.coordinates;
  return `https://google.com/maps/place/${latitude},${longitude}`;
};

export const preferredComment = (
  comment: Record<LanguageCode, string>,
): string | undefined => {
  // const language = navigator.language as LanguageCode;
  return comment["en"] || comment["pl"];
};

export const distanceBetweenPoints = (a: GeoJSON.Point, b: GeoJSON.Point) => {
  return haversine({ geometry: a }, { geometry: b }, { format: "geojson" });
};

export const availabilityColor = (available: number, total: number) => {
  const percent = available / total;
  if (percent > 0.3) return "green";
  else if (percent > 0.1) return "orange";
  else return "red";
};

export const markerColor = (
  available: number,
  total: number,
  status: ParkingLotStatus,
) => {
  return status == ParkingLotStatus.Open
    ? availabilityColor(available, total)
    : parkingLotStatusColor(status);
};
