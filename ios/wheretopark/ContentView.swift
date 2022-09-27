//
//  ContentView.swift
//  Shared
//
//  Created by Grzegorz Barański on 17/05/2022.
//

import SwiftUI
import MapKit
import CoreLocation
import UIKit

struct PrimarySheetView: View {
    @EnvironmentObject var appState: AppState
    @EnvironmentObject var locationManager: LocationManager
    
    var body: some View {
        if appState.fetchFailed {
            Text("Error: \(appState.fetchError?.localizedDescription ?? "")")
            SendFeedback(
                message: {
                    """
                    Hi, I could not open up the application.
                    Error: \(appState.fetchError?.localizedDescription ?? "Unknown")
                    """
                },
                attachment: nil
            )
        } else if appState.parkingLots.isEmpty {
            LoadingView(description: "Loading parking lot data")
        } else {
            ListView()
                .environmentObject(appState)
                .environmentObject(locationManager)
        }
    }
}

struct SecondarySheetView: View {
    @EnvironmentObject var appState: AppState
    
    var body: some View {
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
}

struct ContentView: View {
    @EnvironmentObject var appState: AppState
    @EnvironmentObject var locationManager: LocationManager

    @State var primaryBottomSheetVisible = false
    @State var primaryBottomSheetDetent: UISheetPresentationController.Detent.Identifier? = .compact
    
    @State var secondaryBottomSheetVisible = false
    @State var secondaryBottomSheetDetent: UISheetPresentationController.Detent.Identifier? = .compact
    
    @State var searchText: String = ""
    
    var body: some View {
        NavigationView {
            MapView()
                .edgesIgnoringSafeArea(.all)
                .navigationBarHidden(true)
        }
        .onAppear {
            DispatchQueue.main.asyncAfter(deadline: .now() + 1) {
                primaryBottomSheetVisible = true
            }
        }
        .bottomSheet(
            isPresented: $primaryBottomSheetVisible,
            selectedDetentIdentifier: $primaryBottomSheetDetent,
            isModalInPresentation: true
        ) {
            PrimarySheetView()
                .environmentObject(appState)
                .environmentObject(locationManager)
        }
        .bottomSheet(
            isPresented: $secondaryBottomSheetVisible,
            selectedDetentIdentifier: $secondaryBottomSheetDetent
        ) {
            SecondarySheetView()
                .environmentObject(appState)
        }
        .onChange(of: appState.selectedParkingLotID) { id in
            if id != nil {
                primaryBottomSheetDetent = .small
                secondaryBottomSheetDetent = .compact
                secondaryBottomSheetVisible = true
            } else {
                secondaryBottomSheetVisible = false
                primaryBottomSheetDetent = .compact
            }
        }
        .task({ await appState.fetchParkingLots() })
    }
}

//struct ContentView_Previews: PreviewProvider {
//    static var previews: some View {
//        ContentView()
//    }
//}
