import { Feature } from "./types";

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
    const ofCategoryFeatures = features.filter((f) => OF_CATEGORY_FEATURES.includes(f));
    if (ofCategoryFeatures.length == 0) {
        return "Unknown";
    } else if(ofCategoryFeatures.length == 1) {
        const category = ofCategoryFeatures[0]!;
        if(category == Feature.COVERED) return "Covered";
        else if(category == Feature.UNDERGROUND) return "Underground";
        else if(category == Feature.UNCOVERED) return "Not covered";
    } else {
        if (features.includes(Feature.COVERED) || features.includes(Feature.UNDERGROUND)) {
            if (features.includes(Feature.UNCOVERED)) return "Partially covered";
            else return "Covered & Underground";
        }
    }
    return "Unknown";
}

import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';

dayjs.extend(relativeTime)


export const timeFromNow = (iso: string) => {
    return dayjs(iso).fromNow();
}

export const resourceText = (resource: string) => {
    const url = new URL(resource);
    switch(url.protocol) {
        case "http:":
        case "https:": return url.host;
        case "mailto:": return url.pathname;
        case "tel:": return url.pathname.replaceAll("-", " ");
        default: return "";
    }
};

export const resourceIcon = (resource: string) => {
    const url = new URL(resource);
    switch(url.protocol) {
        case "http:":
        case "https:": return "public";
        case "mailto:": return "mail_outline";
        case "tel:": return "call";
        default: return "error_outline";
    }
};

export const googleMapsLink = (location: GeoJSON.Feature<GeoJSON.Point>) => {
    const [longitude, latitude] = location.geometry.coordinates;
    return `https://google.com/maps/place/${latitude},${longitude}`;

}