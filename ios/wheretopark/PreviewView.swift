//
//  Preview.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 16/11/2022.
//

import SwiftUI
import MapKit


struct PreviewView: View {
    @EnvironmentObject var locationManager: LocationManager
    
    let parkingLot: ParkingLot
    
    var metadata: ParkingLotMetadata { parkingLot.metadata }
    
    var state: ParkingLotState { parkingLot.state }
    
    var body: some View {
        VStack(alignment: .leading) {
            Text(parkingLot.metadata.name)
                .font(.title)
                .fontWeight(.bold)
            Label(metadata.address, systemImage: "map.circle")
                .font(.subheadline)
            Label("\(state.availableSpots[ParkingSpotType.car.rawValue] ?? 0) \(String(localized: "parkingLot.availableSpots"))", systemImage: "parkingsign.circle")
                .font(.subheadline)
            if let userLocation: CLLocation = locationManager.lastLocation {
                let distance = parkingLot.metadata.geometry.distance(from: userLocation)
                let distanceString = MKDistanceFormatter().string(fromDistance: distance)
                Label(distanceString, systemImage: "point.topleft.down.curvedto.point.bottomright.up")
                    .font(.subheadline)
            }
        }
        .padding(20)
        .background(
            RoundedRectangle(cornerRadius: 10)
                .fill(.ultraThickMaterial)
        )
    }
}

struct PreviewView_Previews: PreviewProvider {
    @StateObject static var locationManager = LocationManager()
    
    static var previews: some View {
        ZStack(alignment: .center) {
            Color.green.ignoresSafeArea()
            PreviewView(parkingLot: ParkingLot.example)
                .environmentObject(locationManager)
        }
    }
}

