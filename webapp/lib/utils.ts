export function capitalizeFirstLetter(s: string) {
    return s.charAt(0).toUpperCase() + s.slice(1);
}



import { getDistance, convertDistance  } from 'geolib';
import { Location } from './parkingLot';

export const prettyDistance = (from: Location, to: Location) => {
    console.log({from, to});
    const distance = getDistance(from, to);
    if (distance < 1000) {
        return `${Math.round(distance / 100) * 100} m`;
    } else {
        return `${Math.round(convertDistance(distance, 'km') * 10) / 10} km`;
    }
}