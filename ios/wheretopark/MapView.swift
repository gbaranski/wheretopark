//
//  ParkingLotMap.swift
//  parkflow
//
//  Created by Grzegorz Barański on 18/05/2022.
//

import SwiftUI
import MapKit
import shared

let trackingModeProperties = [
    MKUserTrackingMode.none: ("No user tracking", "location"),
    MKUserTrackingMode.follow: ("Following user", "location.fill"),
    MKUserTrackingMode.followWithHeading: ("Following user with heading", "location.north.line.fill")
]

struct MapView: View {
    @State var userTrackingMode: MKUserTrackingMode = .follow
    @EnvironmentObject var appState: AppState
    
    @Environment(\.colorScheme) private var colorScheme
    
    var body: some View {
            ZStack(alignment: .topTrailing) {
                MapViewRepresentable(
                    userTrackingMode: $userTrackingMode
                )
                VStack {
                    Button(action: {
                        if userTrackingMode == .none {
                            userTrackingMode = .follow
                        } else if userTrackingMode == .follow {
                            userTrackingMode = .followWithHeading
                        } else if userTrackingMode == .followWithHeading {
                            userTrackingMode = .none
                        }
                    }) {
                        let (title, icon) = trackingModeProperties[userTrackingMode]!
                        Label(title, systemImage: icon).labelStyle(.iconOnly)
                    }
                    .foregroundColor(.gray)
                    .padding(.top, 15)
                    
                    Divider()
                        .frame(maxWidth: 44)
                        .background(.gray)
                        .padding(.vertical, 15)
                    
                    Button(action: { Task { await appState.fetchParkingLots() } }) {
                        Label("Refresh", systemImage: "arrow.clockwise").labelStyle(.iconOnly)
                    }
                    .disabled(appState.isPerformingTask)
                    .foregroundColor(.gray)
                    .padding(.bottom, 15)
                }
                .background(colorScheme == .dark ? Color(red: 37 / 255, green: 39 / 255, blue: 42 / 255) : Color.white)
                .cornerRadius(10)
                .padding(.top, 50)
                .padding(20)
            }
    }
}

struct MapViewRepresentable: UIViewRepresentable {
    @EnvironmentObject var appState: AppState
    @Binding var userTrackingMode: MKUserTrackingMode
    @State var didMoveToCurrentLocation = false
    
    let map = MKMapView()
    
    private enum MapDefaults {
        static let latitude = 45.872
        static let longitude = -1.248
        static let zoom = 100.0
        static let zoomWhileSelected = 0.01
    }
    
    func makeUIView(context: Context) -> MKMapView {
        map.delegate = context.coordinator
        let region = MKCoordinateRegion(center: CLLocationCoordinate2D(latitude: 0, longitude: 0), span: MKCoordinateSpan(latitudeDelta: 100.0, longitudeDelta: 100.0))
        self.map.setRegion(region, animated: true)
        self.map.showsTraffic = true
        self.map.showsBuildings = true
        self.map.showsScale = true
        self.map.showsUserLocation = true
        // TODO: Instead of hiding, move it somewhere else to prevent overflowing
        self.map.showsCompass = false
        self.map.userTrackingMode = userTrackingMode
        self.map.directionalLayoutMargins.bottom = 300.0
        return self.map
    }

    func updateAnnotations(view: MKMapView) {
        let (existing, new) = self.appState.parkingLots.reduce(into: ([ParkingLotID : ParkingLot](), [ParkingLotID : ParkingLot]())) {
            let (id, parkingLot) = $1
            let exists = view.annotations
                .compactMap{ $0 as? ParkingLotAnnotation }
                .contains(where: { $0.id == id })
            if exists {
                $0.0[id] = parkingLot
            } else {
                $0.1[id] = parkingLot
            }
        }
        existing.forEach{ id, parkingLot in
            let annotation = view.annotations
                .compactMap{ $0 as? ParkingLotAnnotation }
                .first{ annotation in
                    annotation.id == id
                }
            annotation!.parkingLot = parkingLot
        }
        new.forEach { id, parkingLot in
            let annotation = ParkingLotAnnotation(id: id, parkingLot: parkingLot)
            print("adding annotation for \(id)")
            view.addAnnotation(annotation)
        }
        // remove old parking lots
        view.annotations.filter { annotation in
            if case let annotation as ParkingLotAnnotation = annotation {
                return !self.appState.parkingLots.contains(where: { $0.key == annotation.id })
            } else {
                return false
            }
        }.forEach { annotation in
            view.removeAnnotation(annotation)
        }
    }
    
    func updateUIView(_ view: MKMapView, context: Context) {
        view.delegate = context.coordinator
        if view.userTrackingMode != self.userTrackingMode {
            view.setUserTrackingMode(self.userTrackingMode, animated: true)
        }
        updateAnnotations(view: view)
        if let selectedParkingLotID = self.appState.selectedParkingLotID {
            let annotation = view.annotations
                .compactMap{ $0 as? ParkingLotAnnotation }
                .first{ $0.id == selectedParkingLotID }!
            let userLocationAnnotation = view.annotations.first{ $0 is MKUserLocation }!
            let isAnnotationInVisibleRect = view.annotations(in: view.visibleMapRect).compactMap{ $0 as? ParkingLotAnnotation }.contains{ $0.id == selectedParkingLotID }
            if isAnnotationInVisibleRect {
                view.showAnnotations([annotation], animated: true)
            } else {
                view.showAnnotations([annotation, userLocationAnnotation], animated: true)
            }
            view.selectAnnotation(annotation, animated: true)
        } else {
            view.selectedAnnotations.forEach{ view.deselectAnnotation($0, animated: true) }
        }
    }
    
    func didSelectParkingLot(parkingLotID: ParkingLotID) {
        print("didSelectParkingLot = \(parkingLotID)")
        self.appState.selectedParkingLotID = parkingLotID
    }
    
    func didChangeUserTrackingMode(userTrackingMode: MKUserTrackingMode) {
        self.userTrackingMode = userTrackingMode
    }
    
    func makeCoordinator() -> MapViewCoordinator{
        MapViewCoordinator(self)
    }
    
    class MapViewCoordinator: NSObject, MKMapViewDelegate {
        var mapViewController: MapViewRepresentable
        
        init(_ control: MapViewRepresentable) {
            self.mapViewController = control
        }
        
        func mapView(_ mapView: MKMapView, viewFor annotation: MKAnnotation) -> MKAnnotationView? {
            // TODO: Move it somewhere else
            let color = UIColor(red: 247 / 255, green: 181 / 255, blue: 0 / 255, alpha: 1.0)
            switch annotation {
            case let parkingLotAnnotation as ParkingLotAnnotation:
                var annotationView = mapView.dequeueReusableAnnotationView(withIdentifier: parkingLotAnnotation.id) as? MKMarkerAnnotationView
                if annotationView == nil {
                    annotationView = MKMarkerAnnotationView(annotation: parkingLotAnnotation, reuseIdentifier: parkingLotAnnotation.id)
                } else {
                    annotationView!.annotation = parkingLotAnnotation
                }
                annotationView!.clusteringIdentifier = "parking-lot"
                annotationView!.canShowCallout = false
                annotationView!.animatesWhenAdded = true
                annotationView!.subtitleVisibility = .visible
                annotationView!.markerTintColor = color
                annotationView!.glyphText = String("P")
                return annotationView
            case let cluster as MKClusterAnnotation:
                let markerAnnotationView = MKMarkerAnnotationView()
                let totalAvailableSpots = cluster.memberAnnotations.compactMap{ $0 as? ParkingLotAnnotation }.map{ Int(truncating: $0.parkingLot.state.availableSpots[ParkingSpotType.car] ?? 0) }.reduce(0, +)
                cluster.title = "\(totalAvailableSpots) total available spots"
                cluster.subtitle = nil
                markerAnnotationView.annotation = cluster
                markerAnnotationView.glyphText = String(cluster.memberAnnotations.count)
                markerAnnotationView.markerTintColor = color
                markerAnnotationView.canShowCallout = false
                return markerAnnotationView
            case is MKUserLocation:
                return nil
            default:
                fatalError("received unexpected annotation: \(annotation)")
            }
        }
        
        func mapView(_ mapView: MKMapView, didSelect view: MKAnnotationView) {
            switch view.annotation {
            case let annotation as ParkingLotAnnotation:
                self.mapViewController.didSelectParkingLot(parkingLotID: annotation.id)
            case let annotation as MKClusterAnnotation:
                mapView.showAnnotations(annotation.memberAnnotations, animated: true)
            default:
                break
            }
        }
        
        func mapView(
            _ mapView: MKMapView,
            didChange mode: MKUserTrackingMode,
            animated: Bool
        ) {
            self.mapViewController.didChangeUserTrackingMode(userTrackingMode: mode)
        }
    }
    
}

class ParkingLotAnnotation: NSObject, MKAnnotation {
    var id: ParkingLotID
    var parkingLot: ParkingLot
    var coordinate: CLLocationCoordinate2D {
        self.parkingLot.metadata.location.coordinate
    }
    
    var title: String? {
        self.parkingLot.metadata.name
    }
    
    var subtitle: String? {
        let status = parkingLot.metadata.status()
        if status == .closed || status == .opensSoon {
            return status.name.capitalizingFirstLetter()
        } else {
            return "\(parkingLot.state.availableSpots) available parking spots"
        }
    }
    
    init(id: ParkingLotID, parkingLot: ParkingLot) {
        self.id = id
        self.parkingLot = parkingLot
    }
}

class SubclassedTapGestureRecognizer: UITapGestureRecognizer {
    let parkingLotID: ParkingLotID
    init(target: AnyObject, action: Selector, parkingLotID: ParkingLotID) {
        self.parkingLotID = parkingLotID
        super.init(target: target, action: action)
    }
}
