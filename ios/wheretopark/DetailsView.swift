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
import MessageUI

struct DetailsView: View {
    @EnvironmentObject var appState: AppState
    var onDismiss: (() -> Void)? = nil
    
    let favouriteStore = FavouritesStore()
    var isFavourite: Binding<Bool> {
        Binding {
            self.favouriteStore.exists(id: id)
        } set: { value, tx in
            value ? self.favouriteStore.add(id: id) : self.favouriteStore.remove(id: id)
        }
    }
    var id: ParkingLotID {
        get { appState.selectedParkingLotID! }
    }
    var parkingLot: ParkingLot {
        get { appState.selectedParkingLot.wrappedValue! }
    }
    
    @State private var isSharing = false
    
    var body: some View {
        ScrollView(showsIndicators: false) {
            VStack(alignment: .leading, spacing: 10) {
                if appState.isSelected.wrappedValue {
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
                                    isFavourite.wrappedValue = !isFavourite.wrappedValue
                                }
                            ) {
                                Label(
                                    isFavourite.wrappedValue ? "Remove from favourites" : "Add to favourites",
                                    systemImage: isFavourite.wrappedValue ? "star.fill" : "star"
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
                            Text("\(parkingLot.state.availableSpots[ParkingSpotType.car] ?? 0) cars").fontWeight(.heavy).foregroundColor(.yellow)
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
                            let _ = formatter.locale = Locale(identifier: "en")
                            let lastUpdated = Date(timeIntervalSince1970: TimeInterval(parkingLot.state.lastUpdated.epochSeconds))
                            let lastUpdatedString = formatter.localizedString(for: lastUpdated, relativeTo: Date.now)
                            Text("UPDATED").fontWeight(.black).font(.caption).foregroundColor(.secondary)
                            Text("\(lastUpdatedString)").fontWeight(.bold)
                        }
                        
                    }.frame(maxWidth: .infinity)
                    Text("Pricing").font(.title2).fontWeight(.bold)
                    DetailsPricingView(metadata: parkingLot.metadata)
                    Text("Additional info").font(.title2).fontWeight(.bold)
                    DetailsAdditionalInfoView(id: id, metadata: parkingLot.metadata)
                    DetailsSendFeedbackView(id: id, metadata: parkingLot.metadata, takeSnapshot: { self.environmentObject(appState).snapshot() })
                        .frame(
                            maxWidth: .infinity,
                            maxHeight: .infinity,
                            alignment: .center
                        )
                }
            }
        }.background(SharingViewController(isPresenting: $isSharing) {
            let url = URL(string: UtilitiesKt.getShareURL(id: id))
            let av = UIActivityViewController(activityItems: [url!], applicationActivities: nil)
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
        let mapItem = MKMapItem(placemark: MKPlacemark(coordinate: parkingLot.metadata.location.coordinate, addressDictionary: nil))
        mapItem.name = parkingLot.metadata.name
        mapItem.openInMaps(launchOptions: [MKLaunchOptionsDirectionsModeKey: MKLaunchOptionsDirectionsModeDriving])
    }
}

struct DetailsPricingView: View {
    let metadata: ParkingLotMetadata
    
    var body: some View {
        Group {
            ForEach(Array(metadata.rules.enumerated()), id: \.1) { i, rule in
                if let weekdays: ParkingLotWeekdays = rule.weekdays {
                    Text(weekdays.human()).font(.body).fontWeight(.bold)
                }
                if let hours: ParkingLotHours = rule.hours {
                    Text(hours.human()).font(.caption).fontWeight(.bold)
                }
                ForEach(rule.pricing, id: \.self) { price in
                    let priceString = price.price.formatted(.currency(code: metadata.currency))
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
    }
    
    
}

struct DetailsAdditionalInfoView: View {
    let id: ParkingLotID
    let metadata: ParkingLotMetadata
    
    var body: some View {
        Group {
            VStack(alignment: .leading) {
                Text("Parking lot").foregroundColor(.secondary)
                Text("\(metadata.totalSpots[ParkingSpotType.car] ?? 0) total spaces")
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
            VStack(alignment: .leading) {
                Text("Unique ID").foregroundColor(.secondary)
                Text("\(id)").textSelection(.enabled)
            }
            Divider()
        }
    }
}

struct DetailsSendFeedbackView: View {
    let id: ParkingLotID
    let metadata: ParkingLotMetadata
    let takeSnapshot: () -> UIImage
    
    @State var result: Result<MFMailComposeResult, Error>? = nil
    @State var lastSnapshot: UIImage? = nil
    private var isShowingMailView: Binding<Bool> {
        Binding {
            return lastSnapshot != nil
        } set: { isShowing in
            if (!isShowing) {
                lastSnapshot = nil
            } else {
                fatalError("unexpected isShowing to be true")
            }
        }
        
    }
    
    var body: some View {
        
        Button(action: {
            self.lastSnapshot = takeSnapshot()
        }) {
            Text("Send a feedback")
        }
        .disabled(!MFMailComposeViewController.canSendMail())
        .sheet(isPresented: isShowingMailView) {
            MailView(
                parkingLotID: id,
                image: lastSnapshot!,
                result: self.$result
            )
        }
    }
}


struct DetailsView_Previews: PreviewProvider {
    static var previews: some View {
        DetailsView(
//            id: Binding.constant(ParkingLot.companion.galeriaBaltycka.metadata.location.hash(length: 12)),
//            parkingLot: Binding.constant(ParkingLot.companion.galeriaBaltycka),
        ).padding()
    }
}
