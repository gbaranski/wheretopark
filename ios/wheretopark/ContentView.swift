//
//  ContentView.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 17/05/2022.
//

import SwiftUI
import MapKit
import CoreLocation
import UIKit

struct HorizontalCarousel: View {
    let parkingLots: Array<(key: String, value: ParkingLot)>
    let onSelect: (ParkingLotID) -> Void
    
    var body: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(alignment: .bottom, spacing: 0) {
                ForEach(parkingLots, id: \.key) { id, parkingLot in
                    PreviewView(parkingLot: parkingLot)
                        .padding(.horizontal)
                        .onTapGesture {
                            onSelect(id)
                        }
                }
            }
        }
    }
}

struct VerticalCarousel: View {
    let parkingLots: Array<(key: String, value: ParkingLot)>
    let onSelect: (ParkingLotID) -> Void
    
    var body: some View {
        ScrollView(.vertical, showsIndicators: false) {
            VStack(alignment: .leading, spacing: 0) {
                ForEach(parkingLots, id: \.key) { id, parkingLot in
                    PreviewView(parkingLot: parkingLot)
                        .padding(.vertical)
                        .onTapGesture {
                            onSelect(id)
                        }
                }
            }
        }
    }
}

struct DetailsSheet: View {
    @EnvironmentObject var appState: AppState
    let onDismiss: (() -> Void)
    
    var body: some View {
        if appState.selectedParkingLotID != nil {
            DetailsView(
                id: appState.selectedParkingLotID!,
                parkingLot: appState.parkingLots[appState.selectedParkingLotID!]!,
                onDismiss: onDismiss
            )
            .padding([.top, .horizontal])
        }
    }
}

struct ContentView: View {
    @EnvironmentObject var appState: AppState
    @EnvironmentObject var locationManager: LocationManager
    @Environment(\.verticalSizeClass) var ver: UserInterfaceSizeClass?
    
    @State var searchText: String = ""
    @State var bottomSheetVisible = false
    @State var bottomSheetDetent: UISheetPresentationController.Detent.Identifier? = .compact
    
    func onSelect(id: ParkingLotID) {
        appState.selectedParkingLotID = id
    }
    
    var body: some View {
        let parkingLots = appState.parkingLots.sorted(by: {
            if let userLocation: CLLocation = locationManager.lastLocation {
                return $0.value.metadata.geometry.distance(from: userLocation) < $1.value.metadata.geometry.distance(from: userLocation)
            } else {
                return $0.key > $1.key
            }
        }).filter { id, parkingLot in
            return searchText.isEmpty ? true : parkingLot.metadata.name.lowercased().contains(searchText.lowercased())
        }
        
        ZStack(alignment: ver == .compact ? .leading : .bottom) {
            MapView()
            
            VStack(alignment: .leading) {
                if ver == .compact {
                    VerticalCarousel(parkingLots: parkingLots, onSelect: onSelect)
                } else {
                    HorizontalCarousel(parkingLots: parkingLots, onSelect: onSelect)
                }
                HStack(alignment: .bottom) {
                    Image(systemName: "magnifyingglass")
                    TextField("search", text: $searchText)
                }
                .foregroundColor(Color(UIColor.secondaryLabel))
                .padding()
                .background(
                    RoundedRectangle(cornerRadius: 10)
                        .fill(.ultraThickMaterial)
                )
                .padding(ver == .compact ? .bottom : .all)
            }
            .frame(maxWidth: ver == .compact ? 200 : .infinity)
        }
        .bottomSheet(
            isPresented: $bottomSheetVisible,
            selectedDetentIdentifier: $bottomSheetDetent,
            onDismiss: {
                appState.selectedParkingLotID = nil
            }
        ) {
            DetailsSheet(onDismiss: {
                bottomSheetVisible = false
            }).environmentObject(appState)
        }
        .onChange(of: appState.selectedParkingLotID) { id in
            if id != nil {
                bottomSheetDetent = .compact
                bottomSheetVisible = true
            } else {
                bottomSheetVisible = false
            }
        }
        .alert("Failed to communicate with server", isPresented: $appState.fetchFailed) {
            Button("Retry", role: .cancel) {
                appState.fetchFailed = false
                Task {
                    await appState.fetchParkingLots()
                }
            }
            Button("Exit", role: .destructive) {
                DispatchQueue.main.asyncAfter(deadline: .now()) {
                    UIApplication.shared.perform(#selector(NSXPCConnection.suspend))
                }
            }
        } message: {
            Text(appState.fetchError?.errorDescription ?? "")
        }
        .task({ await appState.fetchParkingLots() })
    }
}

struct ContentView_Previews: PreviewProvider {
    static func initLocationManager() -> LocationManager {
        let locationManager = LocationManager()
        locationManager.lastLocation = CLLocation(latitude:  54.377, longitude: 18.588)
        locationManager.locationStatus = CLAuthorizationStatus.authorizedWhenInUse
        return locationManager
    }
    
    @StateObject static var appState = AppState()
    @StateObject static var locationManager = initLocationManager()
    
    static var previews: some View {
        ContentView()
            .environmentObject(appState)
            .environmentObject(locationManager)
    }
}

