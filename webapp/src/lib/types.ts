export type Coordinate = {
    latitude: number;
    longitude: number;
};

export type Resource = string;

export enum SpotType {
    CAR,
    CAR_DISABLED,
    CAR_ELECTRIC,
    MOTORCYCLE
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
    applies: SpotType[];
    pricing: PricingRule[];

};

export type Metadata = {
    name: string;
    address: string;
    location: GeoJSON.Feature<GeoJSON.Point>;
    resources: Resource[];
    totalSpots: Record<SpotType, number>;
    maxWidth: number | undefined;
    maxHeight: number | undefined;
    features: Feature[];
    paymentMethods: PaymentMethod[];
    comment: Record<LanguageCode, string>;
    currency: string;
    timezone: string;
    rules: Rule[];
};

export type State = {
    // ISO 8601 string
    lastUpdated: string;
    availableSpots: Record<SpotType, number>,
};

export type ParkingLot = {
    metadata: Metadata,
    state: State,
};

export type ID = string;