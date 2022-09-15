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

extension NSNotification.Name {
    static let primarySheetDismiss = NSNotification.Name("com.wheretopark.app.sheet.primary.dismiss")
    static let secondarySheetDismiss = NSNotification.Name("com.wheretopark.app.sheet.secondary.dismiss")
}


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
            selectedDetentIdentifier: $primaryBottomSheetDetent,
            onDismiss: { NotificationCenter.default.post(name: .primarySheetDismiss, object: nil) }
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
            selectedDetentIdentifier: $secondaryBottomSheetDetent,
            onDismiss: { NotificationCenter.default.post(name: .secondarySheetDismiss, object: nil) }
        ) {
            if let parkingLotID = appState.selectedParkingLotID {
                let parkingLot = appState.parkingLots[parkingLotID]!
                    DetailsView(
                        id: parkingLotID,
                        parkingLot: parkingLot,
                        closeAction: {
                            appState.selectedParkingLotID = nil
                        },
                    )
                    .padding([.top, .horizontal])
                }
        }
        .alert(isPresented: $appState.fetchFailed, error: appState.fetchError, actions: {})
        .onChange(of: appState.parkingLots) { _ in
            primaryBottomSheetVisible = true
        }
        .onChange(of: appState.selectedParkingLotID) { id in
            if id != nil {
                primaryBottomSheetVisible = false
                NotificationCenter.default.addObserver(forName: .primarySheetDismiss, object: nil, queue: .main) {_ in
                    secondaryBottomSheetVisible = true
                }
            } else {
                secondaryBottomSheetVisible = false
                NotificationCenter.default.addObserver(forName: .secondarySheetDismiss, object: nil, queue: .main) {_ in
                    primaryBottomSheetVisible = true
                }
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
