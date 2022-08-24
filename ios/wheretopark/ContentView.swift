//
//  ContentView.swift
//  Shared
//
//  Created by Grzegorz Barański on 17/05/2022.
//

import SwiftUI
import BottomSheet
import MapKit
import CoreLocation
import UIKit

extension UISheetPresentationController.Detent {
    open class func small() -> UISheetPresentationController.Detent {
         ._detent(withIdentifier: "small", constant: 80.0)
    }
    
    open class func compact() -> UISheetPresentationController.Detent {
         ._detent(withIdentifier: "compact", constant: 300.0)
    }
}

extension UISheetPresentationController.Detent.Identifier {
    public static let small: UISheetPresentationController.Detent.Identifier = UISheetPresentationController.Detent.Identifier(rawValue: "small")
    
    public static let compact: UISheetPresentationController.Detent.Identifier = UISheetPresentationController.Detent.Identifier(rawValue: "compact")
}

struct ContentView: View {
    @Environment(\.colorScheme) var colorScheme
    @EnvironmentObject var appState: AppState

    @State var primaryBottomSheetVisible = false
    @State var primaryBottomSheetDetent: UISheetPresentationController.Detent.Identifier? = .compact
    
    @State var secondaryBottomSheetDetent: UISheetPresentationController.Detent.Identifier? = .compact
    
    @State var searchText: String = ""
    
    var body: some View {
        NavigationView {
            MapView()
                .edgesIgnoringSafeArea(.all)
                .navigationBarHidden(true)
        }
        .alert(isPresented: $appState.fetchFailed, error: appState.fetchError, actions: {})
        .bottomSheet(
            isPresented: $primaryBottomSheetVisible,
            detents: [.small(), .compact(), .large()],
            largestUndimmedDetentIdentifier: .large,
            prefersGrabberVisible: true,
            selectedDetentIdentifier: $primaryBottomSheetDetent,
            isModalInPresentation: true
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
        }
        .bottomSheet(
            item: $appState.selectedParkingLotID,
            detents: [.small(), .compact(), .large()],
            largestUndimmedDetentIdentifier: .large,
            prefersGrabberVisible: true,
            selectedDetentIdentifier: $secondaryBottomSheetDetent
        ) {
            if let parkingLotID = appState.selectedParkingLotID {
                let parkingLot = appState.parkingLots[parkingLotID]!
                ParkingLotView(
                    parkingLot: parkingLot,
                    closeAction: {
                        appState.selectedParkingLotID = nil
                    }
                ).padding([.top, .horizontal])
            }
        }
        .onChange(of: appState.parkingLots) { _ in
            primaryBottomSheetVisible = true
        }
        .onChange(of: appState.selectedParkingLotID) { newState in
            if newState == nil {
                primaryBottomSheetDetent = .compact
            } else {
                secondaryBottomSheetDetent = .compact
                primaryBottomSheetDetent = .small
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
