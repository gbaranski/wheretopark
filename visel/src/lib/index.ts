import yaml from "js-yaml";

export type SpotType = "CAR"| "CAR_DISABLED";
export const spotTypes: SpotType[] = ['CAR', 'CAR_DISABLED'];

export type ParkingSpot = {
  points: [number, number][];
  type?: SpotType;
}

export const encode = (parkingSpots: ParkingSpot[]): string => {
  return yaml.dump({
    spots: parkingSpots,
  }, {flowLevel: 3});
};

export const decode = (code: string): ParkingSpot[] => {
  const raw = yaml.load(code);
  if (raw === null) throw new Error("null provided");
  if (typeof raw != 'object') throw new Error("object not provided");

  if (!('spots' in raw)) throw new Error("spots not provided");
  if (!Array.isArray(raw.spots)) throw new Error("spots is not an array");
  raw.spots.forEach((spot: unknown) => {
    if (spot === null) throw new Error("spot is null");
    if (typeof spot != 'object') throw new Error("spot is not an object");

    if (!('points' in spot)) throw new Error("points not provided");
    if (!Array.isArray(spot.points)) throw new Error("points is not an array");
    spot.points.forEach((point: unknown) => {
      if (point === null) throw new Error("point is null");
      if (!Array.isArray(point)) throw new Error("point is not an array");
      if (point.length != 2) throw new Error("point is not a 2-tuple");
      if (typeof point[0] != 'number') throw new Error("point[0] is not a number");
      if (typeof point[1] != 'number') throw new Error("point[1] is not a number");
    });
    
    if ('type' in spot) {
      if (typeof spot.type != 'string') throw new Error("type is not a string");
      if (!(spotTypes as string[]).includes(spot.type)) throw new Error("type is not a valid type");
    }
  });


  console.log({raw})
  return raw.spots.map((spot) => ({
    points: spot.points,
    type: spot.type,
  }));
}