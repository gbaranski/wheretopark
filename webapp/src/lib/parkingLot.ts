import OpeningHours, { type argument_hash } from "opening_hours";
import dayjs, { Dayjs } from "dayjs";
import mailtoLink from "mailto-link";

export type Coordinate = {
  latitude: number;
  longitude: number;
};

export class Resource extends URL {
  text(): string {
    switch (this.protocol) {
      case "http:":
      case "https:":
        return this.host;
      case "mailto:":
        return this.pathname;
      case "tel:":
        return this.pathname.replaceAll("-", " ");
      default:
        return "";
    }
  }
}

export class Status {
  private constructor(readonly text: string) {}

  public static open = new Status("Open");
  public static closed = new Status("Closed");
  public static opensSoon = new Status("Opens soon");
  public static closesSoon = new Status("Closes soon");
  
  isOpen() {
    return this.text === Status.open.text;
  }

  isClosed() {
    return this.text === Status.closed.text;
  }
  
  isOpeningSoon() {
    return this.text === Status.opensSoon.text;
  }

  isClosingSoon() {
    return this.text === Status.closesSoon.text;
  }

  
  withComment(comment: string): Status & {comment: string} {
    return Object.assign(Object.create(Status.prototype), this, { comment });
  }
  
  color(): string {
    switch (this.text) {
      case Status.open.text:
        return "#91F5AD";
      case Status.closesSoon.text:
        return "#ec8004";
      case Status.opensSoon.text:
        return "#ec8004";
      case Status.closed.text:
        return "#EF2D56";
      default:
        return "#000000";
    }
  }
}

// export class SpotsMap {
//   constructor(private map: Record<string, number>) {}

//   get(spotType: SpotType): number {
//     return this.map[spotType.codename] || 0;
//   }

//   static fromJSON(json: Record<string, number>): SpotsMap {
//     const array = Object.entries(json)
//       .map(([key, value]): Record<SpotType, number> | undefined => {
//         const spotType = SpotType.fromCodename(key);
//         if (!spotType) {
//           console.error(`unknown spot type: ${key}`);
//           return undefined;
//         } else {
//           return [spotType, value];
//         }
//       })
//       .filter((v) => v != undefined);
//     const map = Object.fromEntries(array);
    
//     return SpotsMap(Object.fromEntries(map));
//   }
// }

// export type AvailableSpots = SpotInfo;
// export type AvailableSpots = SpotInfo;

export class SpotType {
  private constructor(readonly codename: string) {}

  public static car = new SpotType("CAR");
  public static disabledCar = new SpotType("CAR_DISABLED");
  public static electricCar = new SpotType("CAR_ELECTRIC");
  public static motorcycle = new SpotType("MOTORCYCLE");
  public static truck = new SpotType("TRUCK");
  public static bus = new SpotType("BUS");

  public static all = [
    SpotType.car,
    SpotType.disabledCar,
    SpotType.electricCar,
    SpotType.motorcycle,
    SpotType.truck,
    SpotType.bus,
  ];

  static fromCodename(codename: string): SpotType | undefined {
    return this.all.find((sp) => sp.codename == codename);
  }

  icon(): string {
    switch (this) {
      case SpotType.car:
        return "directions_car";
      case SpotType.disabledCar:
        return "accessible";
      case SpotType.electricCar:
        return "electric_bolt";
      case SpotType.motorcycle:
        return "motorcycle";
      case SpotType.truck:
        return "local_shipping";
      case SpotType.bus:
        return "directions_bus";
      default:
        return "error_outline";
    }
  }
}

export enum Feature {
  UNCOVERED = "UNCOVERED",
  COVERED = "COVERED",
  UNDERGROUND = "UNDERGROUND",
  GUARDED = "GUARDED",
}

export const allFeatures = [Feature.UNCOVERED, Feature.COVERED, Feature.UNDERGROUND, Feature.GUARDED];

export enum PaymentMethod {
  CASH = "CASH",
  CARD = "CARD",
  CONTACTLESS = "CONTACTLESS",
  MOBILE = "MOBILE",
}

export type LanguageCode = string;

export type PricingRule = {
  duration: string;
  price: number;
  repeating: boolean;
};

export type Rule = {
  hours: string;
  applies: SpotType[] | undefined;
  pricing: PricingRule[];
};

export type Dimensions = {
  width: number;
  length: number;
  height: number;
};

export class ParkingLot {
  constructor(
    public id: string,
    public name: string,
    public address: string,
    public geometry: GeoJSON.Point,
    public resources: Resource[],
    public totalSpots: Record<string, number>,
    public maxDimensions: Dimensions | undefined,
    public features: Feature[],
    public paymentMethods: PaymentMethod[],
    public comment: Record<LanguageCode, string>,
    public currency: string,
    public timezone: string,
    public rules: Rule[],
    public lastUpdated: Dayjs,
    private availableSpots: Record<string, number>,
  ) {}

  static fromJSON(id: string, json: Record<string, any>): ParkingLot {
    return new ParkingLot(
      id,
      json.metadata.name,
      json.metadata.address,
      json.metadata.geometry,
      json.metadata.resources.map((r: string) => new Resource(r)),
      json.metadata.totalSpots,
      json.metadata.maxDimensions,
      json.metadata.features,
      json.metadata.paymentMethods,
      json.metadata.comment,
      json.metadata.currency,
      json.metadata.timezone,
      json.metadata.rules.map((r: Record<string, any>): Rule => ({
        hours: r.hours,
        applies: r.applies?.map((s: string) => SpotType.fromCodename(s)),
        pricing: r.pricing,
      })),
      dayjs(json.state.lastUpdated),
      json.state.availableSpots,
    );
  }

  category(): string {
    const OF_CATEGORY_FEATURES = [
      Feature.UNCOVERED,
      Feature.COVERED,
      Feature.UNDERGROUND,
    ];

    const ofCategoryFeatures = this.features.filter((f) =>
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
        this.features.includes(Feature.COVERED) ||
        this.features.includes(Feature.UNDERGROUND)
      ) {
        if (this.features.includes(Feature.UNCOVERED)) {
          return "Partially covered";
        } else return "Covered & Underground";
      }
    }
    return "Unknown";
  }

  rulesForDay(
    spotType: SpotType,
    day: number,
  ): (Rule & { humanHours: string })[] {
    const weekday = dayjs().day(day);
    const date = weekday.toDate();
    return this.rules
      .filter((rule) =>
        rule.applies === undefined || rule.applies.includes(spotType)
      )
      .map((rule) => ({ rule, openingHours: new OpeningHours(rule.hours) }))
      .filter(({ rule, openingHours }) => openingHours.getState(date))
      .map(({ rule, openingHours }) => {
        const matchingRuleIndex = openingHours.getMatchingRule(date);
        console.assert(matchingRuleIndex != undefined);
        const humanHours = openingHours.prettifyValue({
          //@ts-ignore
          rule_index: matchingRuleIndex!,
        });
        return { humanHours, ...rule };
      });
  }

  status(spotType: SpotType): Status & {comment: string} {
    const rules = this.rules.filter((rule) =>
      rule.applies === undefined || rule.applies.includes(spotType)
    );
    const allDayRule = rules.find((rule) => rule.hours == "24/7");
    if (allDayRule != null) {
      return Status.open.withComment("24/7");
    }
    const rawOpeningHours = rules.map((rule) => rule.hours).join(", ");
    const openingHours = new OpeningHours(rawOpeningHours);
    const currentDate = dayjs().tz(this.timezone);

    const rawNextChange = openingHours.getNextChange(currentDate.toDate());
    const nextChange = rawNextChange != undefined
      ? dayjs(rawNextChange)
      : undefined;
    const hoursToChange = nextChange
      ? nextChange.diff(currentDate, "hours")
      : undefined;
    if (openingHours.getState(currentDate.toDate())) {
      if (hoursToChange == undefined) {
        return Status.open.withComment("24/7");
      } else if (hoursToChange > 1) {
        return Status.closesSoon.withComment(
          `Closes ${nextChange!.format("HH:mm")}`,
        );
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
        return Status.closesSoon.withComment(matchingRule);
      }
    } else {
      if (hoursToChange == undefined) {
        return Status.closed.withComment("Temporarily closed");
      } else if (hoursToChange < 1) {
        return Status.opensSoon.withComment(
          `Open ${nextChange!.format("HH:mm")}`,
        );
      } else {
        return Status.closed.withComment(
          `Open ${nextChange!.format("dd HH:mm")}`,
        );
      }
    }
  }

  availableSpotsFor(spotType: SpotType): number {
    return this.availableSpots[spotType.codename] || 0;
  }

  totalSpotsFor(spotType: SpotType): number {
    return this.totalSpots[spotType.codename] || 0;
  }

  preferredComment(): string | undefined {
    // const language = navigator.language as LanguageCode;
    return this.comment["en"] || this.comment["pl"];
  }

  availabilityColorFor(spotType: SpotType): string {
    const percent = this.availableSpotsFor(spotType) /
      this.totalSpotsFor(spotType);
    if (percent > 0.3) return "#91F5AD";
    else if (percent > 0.1) return "#ec8004";
    else return "#EF2D56";
  }

  feedbackLink(): URL {
    const url = mailtoLink({
      to: "contact@wheretopark.app",
      subject: `User feedback regarding ${this.name}`,
      body: `Hello, 
  Issue report within wheretopark.app:
  Name: ${this.name}
  Address: ${this.address}
  Date: ${new Date().toISOString()}
  What happened: `,
    });
    return new URL(url);
  }
}

export type ID = string;
