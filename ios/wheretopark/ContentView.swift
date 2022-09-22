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

struct ContentView: View {
    @EnvironmentObject var appState: AppState

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
        .bottomSheet(
            isPresented: $primaryBottomSheetVisible,
            selectedDetentIdentifier: $primaryBottomSheetDetent
        ) {
            VStack {
                HStack {
                    Image(systemName: "magnifyingglass")
                    TextField("Search", text: $searchText)
                }
                .foregroundColor(Color(UIColor.secondaryLabel))
                .padding(.vertical, 8)
                .padding(.horizontal, 5)
                .background(RoundedRectangle(cornerRadius: 10).fill(Color(UIColor.quaternaryLabel)))
                .padding(.top)
                .padding(.horizontal)
                ListView(query: $searchText)
                    .environmentObject(appState)
            }
            .interactiveDismissDisabled(true)
        }
        .bottomSheet(
            isPresented: $secondaryBottomSheetVisible,
            selectedDetentIdentifier: $secondaryBottomSheetDetent
        ) {
            DetailsView(
                onDismiss: {
                    appState.selectedParkingLotID = nil
                }
            )
            .padding([.top, .horizontal])
            .environmentObject(appState)
        }
        .alert(isPresented: $appState.fetchFailed, error: appState.fetchError, actions: {})
        .onChange(of: appState.parkingLots) { _ in
            primaryBottomSheetVisible = true
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
