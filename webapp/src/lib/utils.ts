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