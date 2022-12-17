//
//  ParkingLot.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 30/10/2022.
//

import Foundation
import MapKit
import PhoneNumberKit
import Swift_ISO8601_DurationParser
import DefaultCodable
import CodableGeoJSON
import SwiftUI


typealias ParkingSpotType = String
typealias LanguageCode = String

enum ParkingLotStatus: String, Codable {
    case opensSoon
    case open
    case closesSoon
    case closed
    
    func localizedString() -> String {
        print("rawValue: `\(rawValue)`")
        return NSLocalizedString("parkingLot.status.\(rawValue)", comment: "Status of a parking lot")
    }
    
    func color() -> Color {
        let colors: [Self : Color] = [
            .opensSoon: .yellow,
            .open: .green,
            .closesSoon: .yellow,
            .closed: .red
        ]
        return colors[self]!
    }
}

struct ParkingLotPricingRule: Hashable {
    let duration: DateComponents
    let price: Decimal
    var repeating: Bool = false
    
}

extension ParkingLotPricingRule: Decodable {
    enum CodingKeys: String, CodingKey {
        case duration
        case price
        case repeating
    }
    
    init(from decoder: Decoder) throws {
        let values = try decoder.container(keyedBy: CodingKeys.self)
        let durationString = try values.decode(String.self, forKey: .duration)
        duration = DateComponents.durationFrom8601String(durationString)!
        let priceString = try values.decode(String.self, forKey: .price)
        price = Decimal(string: priceString)!
        repeating = try values.decodeIfPresent(Bool.self, forKey: .repeating) ?? false
    }
}

struct ParkingLotRule: Decodable, Hashable {
    let hours: String
    @Default<Empty>
    var applies: [ParkingSpotType]
    let pricing: [ParkingLotPricingRule]
    
}

struct Dimensions: Decodable {
    static let Empty = Dimensions(width: nil, height: nil, length: nil)
    
    let width: Int?
    let height: Int?
    let length: Int?
}

struct ParkingLotMetadata: Decodable {
    var name: String
    var address: String
    var geometry: GeoJSON.Geometry
    var resources: [URL]
    var totalSpots: [ParkingSpotType : UInt]
    var maxDimensions: Dimensions?
    var features: [String]
    @Default<Empty>
    var paymentMethods: [String]
    @Default<EmptyDictionary>
    var comment: [LanguageCode : String]
    var currency: String
    var timezone: String
    var rules: [ParkingLotRule]
}

struct ParkingLotState: Decodable {
    let lastUpdated: Date
    let availableSpots: [ParkingSpotType : UInt]
}

struct ParkingLot: Decodable {
    let metadata: ParkingLotMetadata
    let state: ParkingLotState
}


extension ParkingLot {
    static let galeriaBaltycka =  ParkingLot(
        metadata: ParkingLotMetadata(
            name: "Galeria Bałtycka",
            address: "ul.Dmowskiego",
            geometry: GeoJSON.Geometry.point(coordinates: PointGeometry.Coordinates(longitude: 18.60024, latitude: 54.38268)),
            resources: [
                URL(string: "mailto:galeria@galeriabaltycka.pl")!,
                URL(string: "tel:+48-58-521-85-52")!,
                URL(string: "https://www.galeriabaltycka.pl/o-centrum/dojazd-parkingi/parkingi/")!
            ],
            totalSpots: [
                "CAR": 1100
            ],
            features: ["COVERED", "UNCOVERED"],
            paymentMethods: Default(wrappedValue: ["CASH", "CONTACTLESS", "CARD"]),
            comment: Default(wrappedValue: [
                "pl": "Na dwóch najwyższych kondygnacjach budynku centrum handlowego oferujemy dwupoziomowy parking i 1100 miejsc postojowych. \n" +
                    "Wjazd do centrum handlowego odbywa się z ronda od strony ulicy Dmowskiego w Gdańsku. \n" +
                    "Komunikację między poziomami parkingowymi a poziomami handlowymi centrum handlowego zapewniają schody ruchome i windy szybkobieżne.\n" +
                    "Prosimy o zachowanie biletu parkingowego i opłacenie należności za postój w kasie automatycznej, znajdującej się przy wyjściu z parkingu.",
            ]),
            currency: "PLN",
            timezone: "Europe/Warsaw",
            rules: [
                ParkingLotRule(
                    hours: "Mo-Sa 08:00-22:00; Su 09:00-21:00",
                    applies: Default(wrappedValue: []),
                    pricing: [
                        ParkingLotPricingRule(duration: DateComponents(hour: 1), price: 0),
                        ParkingLotPricingRule(duration: DateComponents(hour: 2), price: 2),
                        ParkingLotPricingRule(duration: DateComponents(hour: 3), price: 5),
                        ParkingLotPricingRule(duration: DateComponents(hour: 1), price: 4, repeating: true),
                        ParkingLotPricingRule(duration: DateComponents(day: 1), price: 25),
                    ]
                ),
            ]
        ),
        state: ParkingLotState(
            lastUpdated: Calendar.current.date(byAdding: .second, value: -10, to: Date.now)!,
            availableSpots: [
                "CAR": 123
            ]
        )
    )
}
