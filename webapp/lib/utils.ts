export function capitalizeFirstLetter(s: string) {
    return s.charAt(0).toUpperCase() + s.slice(1);
}

import { getDistance, convertDistance  } from 'geolib';
import { Coordinate } from './types';

export const prettyDistance = (from: Coordinate, to: Coordinate) => {
    console.log({from, to});
    const distance = getDistance(from, to);
    if (distance < 1000) {
        return `${Math.round(distance / 100) * 100} m`;
    } else {
        return `${Math.round(convertDistance(distance, 'km') * 10) / 10} km`;
    }
}