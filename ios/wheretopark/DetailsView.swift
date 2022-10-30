//
//  DetailsView.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 30/10/2022.
//  Created by Grzegorz Barański on 19/05/2022.
//

import SwiftUI
import MapKit
import PhoneNumberKit
import MessageUI

struct DetailsView: View {
    let id: ParkingLotID
    let parkingLot: ParkingLot
    var onDismiss: (() -> Void)? = nil
    
    @State private var isSharing = false
    @ObservedObject var favouriteManager: FavouriteManager
    var isFavourite: Bool {
        get { favouriteManager.isFavourite }
    }
    
    init(id: ParkingLotID, parkingLot: ParkingLot, onDismiss: (() -> Void)? = nil) {
        self.id = id
        self.parkingLot = parkingLot
        self.onDismiss = onDismiss
        self.favouriteManager = FavouriteManager(id: id)
    }
    
    var body: some View {
        ScrollView(showsIndicators: false) {
            VStack(alignment: .leading, spacing: 10) {
                HStack {
                    Text(parkingLot.metadata.name).font(.title).fontWeight(.bold).multilineTextAlignment(.leading)
                    if let onDismiss = onDismiss {
                        Spacer()
                        Button(action: onDismiss) {
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
                }
                HStack {
                    Button(action: navigate) {
                        Label("Navigate", systemImage: "car.fill").frame(maxWidth: .infinity)
                    }
                    .controlSize(.large)
                    .buttonStyle(.borderedProminent)
                    .frame(maxWidth: .infinity)
                    
                    Menu {
                        Button(
                            action: {
                                isFavourite ? favouriteManager.remove() : favouriteManager.add()
                            }
                        ) {
                            Label(
                                isFavourite ? "Remove from favourites" : "Add to favourites",
                                systemImage: isFavourite ? "star.fill" : "star"
                            )
                        }
                        Button(action: {
                            isSharing = true
                        }) {
                          Label("Share", systemImage: "square.and.arrow.up")
                        }
                    } label: {
                        Button(action: {}) {
                            Label("More", systemImage: "ellipsis")
                        }
                        .controlSize(.large)
                        .buttonStyle(.bordered)
                    }
                }
                HStack{
                    VStack(alignment: .leading) {
                        Text("AVAILABILITY").fontWeight(.black).font(.caption).foregroundColor(.secondary)
                        Group {
                            let availableSpots = parkingLot.state.availableSpots["CAR"] ?? 0
                            let totalSpots = parkingLot.metadata.totalSpots["CAR"] ?? 0
                            let color = availabilityColor(available: availableSpots, total: totalSpots)
                            Text("\(availableSpots)").fontWeight(.heavy).foregroundColor(color) +
                            Text(" / \(totalSpots) cars").fontWeight(.heavy).foregroundColor(color).font(.caption)
                        }
                    }
                    Divider()
                    VStack(alignment: .leading) {
                        Text("HOURS").fontWeight(.black).font(.caption).foregroundColor(.secondary)
                        // TODO: Use real status
                        let status = ParkingLotStatus.closed
                        switch status {
                        case .opensSoon:
                            Text("Opens soon").fontWeight(.heavy).foregroundColor(.yellow)
                        case .open:
                            Text("Open").fontWeight(.heavy).foregroundColor(.green)
                        case .closesSoon:
                            Text("Closes soon").fontWeight(.heavy).foregroundColor(.yellow)
                        case .closed:
                            Text("Closed").fontWeight(.heavy).foregroundColor(.red)
//                        default:
//                            fatalError("unknown status \(status)")
                        }
                    }
                    Divider()
                    VStack(alignment: .leading) {
                        let formatter = RelativeDateTimeFormatter()
                        let _ = formatter.locale = Locale(identifier: "en")
                        let lastUpdatedString = formatter.localizedString(for: parkingLot.state.lastUpdated, relativeTo: Date.now)
                        Text("UPDATED").fontWeight(.black).font(.caption).foregroundColor(.secondary)
                        Text("\(lastUpdatedString)").fontWeight(.bold)
                    }
                    
                }.frame(maxWidth: .infinity)
                Text("Pricing").font(.title2).fontWeight(.bold)
                DetailsRulesView(metadata: parkingLot.metadata)
                Text("Additional info").font(.title2).fontWeight(.bold)
                DetailsAdditionalInfo(id: id, metadata: parkingLot.metadata)
                DetailsSendFeedbackView(id: id, metadata: parkingLot.metadata, takeSnapshot: { self.snapshot() })
                    .frame(
                        maxWidth: .infinity,
                        maxHeight: .infinity,
                        alignment: .center
                    )
            }
        }.background(SharingViewController(isPresenting: $isSharing) {
            let url = getShareURL(id: id)
            let av = UIActivityViewController(activityItems: [url], applicationActivities: nil)
            // for iPad
            if UIDevice.current.userInterfaceIdiom == .pad {
               av.popoverPresentationController?.sourceView = UIView()
            }
            av.completionWithItemsHandler = { _, _, _, _ in
                   isSharing = false // required for re-open !!!
               }
            return av
        })
    }
    
    func navigate() {
        let mapItem = MKMapItem(placemark: MKPlacemark(coordinate: parkingLot.metadata.geometry.location!.coordinate, addressDictionary: nil))
        mapItem.name = parkingLot.metadata.name
        mapItem.openInMaps(launchOptions: [MKLaunchOptionsDirectionsModeKey: MKLaunchOptionsDirectionsModeDriving])
    }
}

struct DetailsRulesView: View {
    let metadata: ParkingLotMetadata
    
    var body: some View {
        Group {
            ForEach(Array(metadata.rules.enumerated()), id: \.1) { i, rule in
                DetailsRuleView(rule: rule, currency: metadata.currency)
            }
        }
    }
}

struct DetailsRuleView: View {
    let rule: ParkingLotRule
    let currency: String
    
    var body: some View {
        HStack(alignment: .center) {
            VStack(alignment: .leading) {
                ForEach(Array(rule.expandedHours.enumerated()), id: \.1) { i, hours in
                    Text(hours.trimmingCharacters(in: .whitespacesAndNewlines)).font(.body).fontWeight(.bold)
                }
            }
            HStack {
//                ForEach(Array(rule.applies.enumerated()), id: \.1) { i, spotType in
//                    switch(spotType) {
//                    case .car:
//                        Image(systemName: "car.fill")
//                    case .carDisabled:
//                        Text("♿️")
//                    case .carElectric:
//                        Image(systemName: "bolt.car.fill")
//                    case .motorcycle:
//                        Image(systemName: "bicycle")
//                    default:
//                        Image(systemName: "questionmark.diamond")
//                    }
//                    if i != (rule.applies.count ?? 0) - 1 {
//                        Divider()
//                    }
//                }
            }.frame(
                maxWidth: .infinity,
                alignment: .topTrailing
            )
        }
        DetailsRulePricingView(rule: rule, currency: currency)
    }
}

struct DetailsRulePricingView: View {
    let rule: ParkingLotRule
    let currency: String
    
    var body: some View {
        ForEach(rule.pricing, id: \.self) { price in
            let priceString = price.price.formatted(.currency(code: currency))
            HStack {
                if price.repeating {
                    Text("Each additional \(price.durationString)")
                } else {
                    Text("\(price.durationString)")
                }
                Spacer()
                if price.price == 0 {
                    Text("Free").bold()
                } else {
                    Text(priceString)
                }
            }
            Divider()
        }
    }
    
}

struct DetailsAdditionalInfoField<ContentView: View>: View {
    let name: String
    @ViewBuilder let contentView: () -> ContentView
    
    var body: some View {
        VStack(alignment: .leading) {
            Text(name).foregroundColor(.secondary)
            contentView()
            Divider()
        }
    }
}

struct DetailsAdditionalInfo: View {
    let id: ParkingLotID
    let metadata: ParkingLotMetadata
    
    var body: some View {
        Group {
            DetailsAdditionalInfoField(name: "Parking lot") {
                Text("\(metadata.totalSpots["CAR"] ?? 0) total spaces")
            }
            DetailsAdditionalInfoField(name: "Address") {
                Text("\(metadata.address)")
            }
            DetailsAdditionalInfoField(name: "Coordinates") {
                Text("\(metadata.geometry.location!.coordinate.latitude), \(metadata.geometry.location!.coordinate.longitude)")
            }
            ForEach(metadata.resources, id: \.self) { resource in
//                DetailsAdditionalInfoField(name: resource.label()) {
//                    Link(
//                        "\(resource.components.host ?? "")\(resource.components.scheme == "tel" ? resource.components.path.replacingOccurrences(of: "-", with: " ") : resource.components.path)",
//                        destination: resource.components.url!).truncationMode(.tail).lineLimit(1)
//                }
            }
            if let comment = metadata.commentForLocale {
                DetailsAdditionalInfoField(name: "Comment") {
                    Text(comment)
                }
            }
            DetailsAdditionalInfoField(name: "Unique ID") {
                Text("\(id)").textSelection(.enabled)
            }
        }
    }
}

struct DetailsSendFeedbackView: View {
    let id: ParkingLotID
    let metadata: ParkingLotMetadata
    let takeSnapshot: () -> UIImage
    
    var body: some View {
        SendFeedback(
            message: {
                """
                <p>Hi, I've got a problem with parking lot of ID: \(id)</p>
                <br/>
                <p>My problem is: (describe your problem here)</p>
                """
            },
            attachment: {
                takeSnapshot()
            }
        )
    }
}


struct DetailsView_Previews: PreviewProvider {
    static var previews: some View {
        DetailsView(
            id: "u3tjrk061424",
            parkingLot: ParkingLot.galeriaBaltycka
        ).padding([.horizontal])
    }
}
