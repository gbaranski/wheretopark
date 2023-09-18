import yaml from "js-yaml";

export type Polygon = { x: number; y: number }[];

type Code = {
  spots: {
    points: [number, number][];
  }[];
};

export const generateCode = (
  polygons: Polygon[],
): string => {
  const code: Code = {
    spots: polygons.map((poly) => ({
      points: poly.map(({ x, y }) => [x, y]),
    })),
  };
  return yaml.dump({
    spots: code.spots,
    
  }, {flowLevel: 3});
};

export const parseCode = (code: string): Polygon[] => {
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
  });


  console.log({raw})
  return raw.spots.map((spot) => spot.points.map(([x, y]: [number, number]) => ({ x, y })));
}