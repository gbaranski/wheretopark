export type ID = string;

export type ParkingLot = {
    id: ID;
    metadata: Metadata;
    state: State;
}

export type Metadata = {
    name: string;
    address: string,
    location: Location,
    emails: string[],
    phoneNumbers: string[],
    websites: URL[],
    totalSpots: number | null
    features: Feature[],
    currency: string,
    rules: Rule[],
}

export type State = {
    lastUpdated: string;
    availableSpots: string;
}

export type Location = {
    latitude: number;
    longitude: number;
}

export enum Feature {
    Uncovered,
    Covered,
    Underground
}

export enum Weekday {
    Monday,
    Tuesday,
    Wednesday,
    Thursday,
    Friday,
    Saturday,
    Sunday,
}

export type Weekdays = {
    start: Weekday,
    end: Weekday,
}

export type PricingRule = {
    duration: string;
    price: number;
    repeating: boolean;
}

export type Hours = {
    start: string,
    end: string,
}

export type Rule = {
    weekdays: Weekdays;
    hours: Hours | null;
    pricing: PricingRule[];
}
