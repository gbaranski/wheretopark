//
//  ParkingLotView.swift
//  parkflow
//
//  Created by Grzegorz Barański on 19/05/2022.
//

import SwiftUI
import MapKit
import PhoneNumberKit
import shared

struct ParkingLotView: View {
    let parkingLot: ParkingLot
    var closeAction: (() -> Void)? = nil
    
    var body: some View {
        ScrollView(showsIndicators: false) {
            VStack(alignment: .leading, spacing: 10) {
                if let closeAction = closeAction {
                    HStack {
                        Text(parkingLot.metadata.name).font(.title).fontWeight(.bold).multilineTextAlignment(.leading)
                        Spacer()
                        Button(action: closeAction) {
                            ZStack {
                                Circle()
                                    .frame(width: 30, height: 30)
                                    .foregroundColor(Color(.secondarySystemFill))
                                Image(systemName: "xmark")
                                    .font(Font.body.weight(.bold))
                                    .foregroundColor(.secondary)
                                    .imageScale(.medium)
                                    .frame(width: 44, height: 44)
                            }
                        }
                            
                    }
                } else {
                    Text(parkingLot.metadata.name).font(.title).fontWeight(.bold).multilineTextAlignment(.leading)
                }
                HStack {
                    Button(action: navigate) {
                        Label("Navigate", systemImage: "car.fill").frame(maxWidth: .infinity)
                    }
                    .controlSize(.large)
                    .buttonStyle(.borderedProminent)
                    .frame(maxWidth: .infinity)
                    
                    Menu {
                        Button(action: {}) {
                            Label("Add to favourites", systemImage: "star")
                        }
                    } label: {
                        Button(action: addToFavourites) {
                            Label("More", systemImage: "ellipsis")
                        }
                        .controlSize(.large)
                        .buttonStyle(.bordered)
                    }
                }
                HStack{
                    VStack(alignment: .leading) {
                        Text("AVAILABILITY").fontWeight(.black).font(.caption).foregroundColor(.secondary)
                        Text("\(parkingLot.state.availableSpots) cars").fontWeight(.heavy).foregroundColor(.yellow)
                    }
                    Divider()
                    VStack(alignment: .leading) {
                        Text("HOURS").fontWeight(.black).font(.caption).foregroundColor(.secondary)
                        let status = parkingLot.metadata.status()
                        switch status {
                        case .opensSoon:
                            Text("Opens soon").fontWeight(.heavy).foregroundColor(.yellow)
                        case .open:
                            Text("Open").fontWeight(.heavy).foregroundColor(.green)
                        case .closesSoon:
                            Text("Closes soon").fontWeight(.heavy).foregroundColor(.yellow)
                        case .closed:
                            Text("Closed").fontWeight(.heavy).foregroundColor(.red)
                        default:
                            fatalError("unknown status \(status)")
                        }
                    }
                    Divider()
                    VStack(alignment: .leading) {
                        let formatter = RelativeDateTimeFormatter()
                        let lastUpdated = Date(timeIntervalSince1970: TimeInterval(parkingLot.state.lastUpdated.epochSeconds))
                        let lastUpdatedString = formatter.localizedString(for: lastUpdated, relativeTo: Date.now)
                        Text("UPDATED").fontWeight(.black).font(.caption).foregroundColor(.secondary)
                        Text("\(lastUpdatedString)").fontWeight(.bold)
                    }
                    
                }.frame(maxWidth: .infinity)
                Text("Pricing").font(.title2).fontWeight(.bold)
                ParkingLotPricingView(metadata: parkingLot.metadata)
                Text("Additional info").font(.title2).fontWeight(.bold)
                ParkingLotAdditionalInfoView(metadata: parkingLot.metadata)
            }
        }
    }
    
    func addToFavourites() {
        
    }
    
    func navigate() {
        let mapItem = MKMapItem(placemark: MKPlacemark(coordinate: parkingLot.metadata.location.coordinate, addressDictionary: nil))
        mapItem.name = parkingLot.metadata.name
        mapItem.openInMaps(launchOptions: [MKLaunchOptionsDirectionsModeKey: MKLaunchOptionsDirectionsModeDriving])
    }
}

struct ParkingLotPricingView: View {
    let metadata: ParkingLotMetadata
    
    var body: some View {
        Group {
            ForEach(Array(metadata.rules.enumerated()), id: \.1) { i, rule in
//                if let weekdays: ParkingLotWeekdays = rule.weekdays {
//                    Text(weekdays.intervalString).font(.body).fontWeight(.bold)
//                }
//                if let hours: ParkingLotHours = rule.hours {
//                    Text(hours.intervalString).font(.caption).fontWeight(.bold)
//                }
//                ForEach(rule.pricing, id: \.self) { price in
//                    let priceString = price.price.formatted(.currency(code: metadata.currency))
//                    HStack {
//                        if price.repeating {
//                            Text("Each additional \(price.durationString)")
//                        } else {
//                            Text("\(price.durationString)")
//                        }
//                        Spacer()
//                        if price.price == 0 {
//                            Text("Free").bold()
//                        } else {
//                            Text(priceString)
//                        }
//                    }
//                    Divider()
//                }
            }
        }
    }
    
    
}

struct ParkingLotAdditionalInfoView: View {
    let metadata: ParkingLotMetadata
    
    var body: some View {
        Group {
            VStack(alignment: .leading) {
                Text("Parking lot").foregroundColor(.secondary)
                Text("\(metadata.totalSpots) total spaces")
            }
            Divider()
            VStack(alignment: .leading) {
                Text("Address").foregroundColor(.secondary)
                Text("\(metadata.address)")
            }
            Divider()
            VStack(alignment: .leading) {
                Text("Coordinates").foregroundColor(.secondary)
                Text("\(metadata.location.latitude), \(metadata.location.longitude)")
            }
            ForEach(metadata.resources, id: \.self) { resource in
                Divider()
                VStack(alignment: .leading) {
                    Text(resource.label()).foregroundColor(.secondary)
                    Link(
                        "\(resource.components.host ?? "")\(resource.components.scheme == "tel" ? resource.components.path.replacingOccurrences(of: "-", with: " ") : resource.components.path)",
                        destination: resource.components.url!).truncationMode(.tail).lineLimit(1)
                }
            }
            Divider()
        }
    }
}


struct ParkingLotView_Previews: PreviewProvider {
    static var previews: some View {
        ParkingLotView(parkingLot: ParkingLot.companion.galeriaBaltycka, closeAction: {}).padding()
    }
}
