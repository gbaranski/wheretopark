//
//  ParkingLot.swift
//  iosApp
//
//  Created by Grzegorz Barański on 23/08/2022.
//  Copyright © 2022 orgName. All rights reserved.
//

import Foundation
import shared
import CoreLocation
import PhoneNumberKit

extension Coordinate {
    var coordinate: CLLocationCoordinate2D {
        CLLocationCoordinate2D(latitude: latitude, longitude: longitude)
    }
    
    func distance(from: CLLocation) -> CLLocationDistance {
        let pointLocation = CLLocation(
            latitude: self.latitude,
            longitude: self.longitude
        )
        return from.distance(from: pointLocation)
    }
}

extension ParkingLotResource {
    var components: URLComponents {
        URLComponents(string: self.url)!
    }
}


extension ParkingLotPricingRule {
    var durationString: String {
        let duration: Duration = .nanoseconds(self.duration)
        return duration.formatted(
            .units(allowed: [.hours, .minutes, .seconds, .milliseconds],
                   width: .wide)
            .locale(Locale(identifier: "en"))
        )
    }
}
