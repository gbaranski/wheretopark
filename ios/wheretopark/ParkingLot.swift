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

extension ParkingLotLocation {
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
