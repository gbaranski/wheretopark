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

struct Carousel: View {
    @EnvironmentObject var appState: AppState
    @EnvironmentObject var locationManager: LocationManager
    @Binding var searchText: String
    
    var body: some View {
        let processedParkingLots = appState.parkingLots.sorted(by: {
            if let userLocation: CLLocation = locationManager.lastLocation {
                return $0.value.metadata.geometry.distance(from: userLocation) < $1.value.metadata.geometry.distance(from: userLocation)
            } else {
                return $0.key > $1.key
            }
        }).filter { id, parkingLot in
            return searchText.isEmpty ? true : parkingLot.metadata.name.lowercased().contains(searchText.lowercased())
        }
        
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(alignment: .bottom, spacing: 0) {
                ForEach(processedParkingLots, id: \.key) { id, parkingLot in
                    PreviewView(parkingLot: parkingLot)
                        .padding(.horizontal)
                        .onTapGesture {
                            appState.selectedParkingLotID = id
                        }
                }
            }
        }
    }
    
    
}

struct ContentView: View {
    @EnvironmentObject var appState: AppState
    @EnvironmentObject var locationManager: LocationManager
    
    @State var searchText: String = ""
    @State var bottomSheetVisible = false
    @State var bottomSheetDetent: UISheetPresentationController.Detent.Identifier? = .compact
    
    var body: some View {
        
        ZStack(alignment: .bottom) {
            MapView()
                .edgesIgnoringSafeArea(.all)
                .navigationBarHidden(true)
            
            VStack {
                Carousel(searchText: $searchText)
                HStack {
                    Image(systemName: "magnifyingglass")
                    TextField("search", text: $searchText)
                }
                .foregroundColor(Color(UIColor.secondaryLabel))
                .padding()
//                .padding(.vertical, 8)
//                .padding(.horizontal, 5)
                .background(
                    RoundedRectangle(cornerRadius: 10)
                        .fill(.ultraThickMaterial)
                )
                .padding()
            }
        }
        .bottomSheet(
            isPresented: $bottomSheetVisible,
            selectedDetentIdentifier: $bottomSheetDetent
        ) {
            if appState.selectedParkingLotID != nil && appState.selectedParkingLot.wrappedValue != nil {
                DetailsView(
                    id: appState.selectedParkingLotID!,
                    parkingLot: appState.selectedParkingLot.wrappedValue!,
                    onDismiss: {
                        appState.selectedParkingLotID = nil
                    }
                )
                .padding([.top, .horizontal])
                .environmentObject(appState)
            }
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

