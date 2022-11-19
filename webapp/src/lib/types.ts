export type Coordinate = {
    latitude: number;
    longitude: number;
};

export type Resource = string;

export enum ParkingLotStatus {
    Open = "Open",
    Closed = "Closed",
    OpensSoon = "Opens soon",
    ClosesSoon = "Closes soon",
}

export enum SpotType {
    CAR = "CAR",
    CAR_DISABLED = "CAR_DISABLED",
    CAR_ELECTRIC = "CAR_ELECTRIC",
    MOTORCYCLE = "MOTORCYCLE",
}
    
export enum Feature {
    UNCOVERED,
    COVERED,
    UNDERGROUND,
    GUARDED,
}

export enum PaymentMethod {
    CASH,
    CARD,
    CONTACTLESS,
    MOBILE,
}

export type LanguageCode = string;

export type PricingRule = {
    duration: string;
    price: number;
    repeating: boolean;
}

export type Rule = {
    hours: string;
    applies: SpotType[] | undefined;
    pricing: PricingRule[];

};

export type Metadata = {
    name: string;
    address: string;
    geometry: GeoJSON.Point;
    resources: Resource[];
    totalSpots: Record<string, number>;
    maxWidth: number | undefined;
    maxHeight: number | undefined;
    features: string[];
    paymentMethods: string[];
    comment?: Record<LanguageCode, string>;
    currency: string;
    timezone: string;
    rules: Rule[];
};

export type State = {
    // ISO 8601 string
    lastUpdated: string;
    availableSpots: Record<string, number>,
};

export type ParkingLot = {
    metadata: Metadata,
    state: State,
};

export type ID = string;