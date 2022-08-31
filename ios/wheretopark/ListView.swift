//
//  ParkingLotList.swift
//  parkflow
//
//  Created by Grzegorz Barański on 26/05/2022.
//

import SwiftUI
import CoreLocation
import MapKit
import shared

struct ListView: View {
    @Binding var query: String
    @EnvironmentObject var appState: AppState
    
    var body: some View {
        let processedParkingLots = appState.parkingLots.sorted(by: {
            if let userLocation: CLLocation = appState.locationManager.lastLocation {
                return $0.value.metadata.location.distance(from: userLocation) < $1.value.metadata.location.distance(from: userLocation)
            } else {
                return $0.key > $1.key
            }
        }).filter { id, parkingLot in
            return query.isEmpty ? true : parkingLot.metadata.name.lowercased().contains(query.lowercased())
        }
        List {
            ForEach(processedParkingLots, id: \.key) { id, parkingLot in
                VStack(alignment: .leading, spacing: 3) {
                    Text(parkingLot.metadata.name).foregroundColor(.primary).font(.headline)
                    Label("\(parkingLot.state.availableSpots[ParkingSpotType.car] ?? 0) available parking spots", systemImage: "parkingsign.circle")
                        .foregroundColor(.secondary)
                        .font(.subheadline)
                    if let userLocation: CLLocation = appState.locationManager.lastLocation {
                        let distance = parkingLot.metadata.location.distance(from: userLocation)
                        let distanceString = MKDistanceFormatter().string(fromDistance: distance)
                        Label("\(distanceString) away", systemImage: "point.topleft.down.curvedto.point.bottomright.up")
                            .foregroundColor(.secondary)
                            .font(.subheadline)
                    }
                }.onTapGesture {
                    appState.selectedParkingLotID = id
                }
            }
            
        }
    }
}
