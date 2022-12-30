//
//  Utilities.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 21/06/2022.
//

import Foundation
import MapKit
import SwiftUI
import CodableGeoJSON

extension GeoJSON.Geometry {
    var location: CLLocation? {
        if case GeoJSON.Geometry.point(let coordinates) = self {
            return CLLocation(latitude: coordinates.latitude, longitude: coordinates.longitude)
        } else {
            return nil
        }
    }
    
    func distance(from: CLLocation) -> CLLocationDistance {
        return self.location!.distance(from: from)
    }
}


extension View {
    func snapshot() -> UIImage {
        let controller = UIHostingController(rootView: self)
        let view = controller.view

        let targetSize = controller.view.intrinsicContentSize
        view?.bounds = CGRect(origin: .zero, size: targetSize)
        view?.backgroundColor = .clear

        let renderer = UIGraphicsImageRenderer(size: targetSize)

        return renderer.image { _ in
            view?.drawHierarchy(in: controller.view.bounds, afterScreenUpdates: true)
        }
    }
}

struct SharingViewController: UIViewControllerRepresentable {
    @Binding var isPresenting: Bool
    var content: () -> UIViewController
    
    func makeUIViewController(context: Context) -> UIViewController {
        UIViewController()
    }
    
    func updateUIViewController(_ uiViewController: UIViewController, context: Context) {
        if isPresenting {
            uiViewController.present(content(), animated: true, completion: nil)
        }
    }
}

extension String {
    func capitalizingFirstLetter() -> String {
      return prefix(1).uppercased() + self.lowercased().dropFirst()
    }

    mutating func capitalizeFirstLetter() {
      self = self.capitalizingFirstLetter()
    }
}

func availabilityColor(available: UInt, total: UInt) -> Color {
    let percent = Double(available) / Double(total)
    if percent > 0.3 {
        return .green
    } else if percent > 0.1 {
        return .yellow
    } else {
        return .red
    }
}


let BASE_WEBAPP_URL: URL = URL(string: "https://web.wheretopark.app")!

func getShareURL(id: ParkingLotID) -> URL {
    return BASE_WEBAPP_URL.appending(path: "/parking-lot/\(id)")
}


extension ParkingLotRule {
    var expandedHours: [String] {
        return self.hours.components(separatedBy: ";")
    }
    
}

extension ParkingLotPricingRule {
    var durationString: String {
        let durationFormatter = DateComponentsFormatter()
        durationFormatter.unitsStyle = .full
        return durationFormatter.string(from: self.duration)!
    }
}


import OpeningHours

extension ParkingLotMetadata {
    var commentForLocale: String? {
        let languageCode = Locale.current.language.languageCode?.identifier ?? "en"
        let comment = comment[languageCode] ?? comment ["en"]
        return comment
    }
    
    func status() -> ParkingLotStatus {
        let rawOpeningHours = rules.map{ $0.hours }.joined(separator: ", ")
        let openingHours = OpeningHours(rawOpeningHours)
        let dateTime = now_at(timezone)
        if openingHours.is_open(dateTime) {
            return .open
        } else {
            return .closed
        }
    }
}

extension ParkingLot {
    func tintColor() -> Color {
        let availabilityColor = availabilityColor(available: state.availableSpots[ParkingSpotType.car.rawValue]!, total: metadata.totalSpots[ParkingSpotType.car.rawValue]!)
        let status = metadata.status()
        let tintColor = status == .closed ? .red : status == .closesSoon ? .yellow : availabilityColor
        return tintColor
    }
    
}

extension URL {
    var label: LocalizedStringKey {
        let labels = [
            "http": "url.website",
            "https": "url.website",
            "tel": "url.phoneNumber",
            "mailto": "url.emailAddress",
        ];
        return LocalizedStringKey(labels[scheme!] ?? "url.unrecognized");
    }
    
    var human: String {
        if (scheme == "http" || scheme == "https") {
            return "\(host!)\(path)";
        } else if (scheme == "mailto") {
            return String(absoluteString.dropFirst(7));
        } else if (scheme == "tel") {
            return String(absoluteString.dropFirst(4)).replacingOccurrences(of: "-", with: " ");
        } else {
            return absoluteString;
        }
    }
}
