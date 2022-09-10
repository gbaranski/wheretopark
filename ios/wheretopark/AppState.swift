//
//  AppState.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 11/07/2022.
//

import Foundation
import shared
import SwiftUI
import Contacts

typealias ParkingLotID = String

@MainActor class AppState: ObservableObject {
    let environment = AppEnvironment()
    var authorizationClient: AuthorizationClient {
        get {
            AuthorizationClient(baseURL: environment.authorizationURL.absoluteString, clientID: environment.clientID, clientSecret: environment.clientSecret)
        }
    }
    var storekeeperClient: StorekeeperClient {
        get {
            StorekeeperClient(
                baseURL: environment.storekeeperURL.absoluteString,
                authorizationClient: authorizationClient,
                accessScope: Set([AccessType.readmetadata, AccessType.readstate])
            )
        }
    }
    @Published private(set) var locationManager = LocationManager()
    @Published private(set) var isPerformingTask = false
    @Published private(set) var fetchError: FetchError? = nil
    @Published var fetchFailed = false
    
    @Published private(set) var parkingLots = [ParkingLotID : ParkingLot]()
    @Published var selectedParkingLotID: ParkingLotID? = nil
    
    
    struct FetchError: LocalizedError {
        let errorDescription: String?
        // TODO: Make this so it's displayed correctly on the Alert
        let recoverySuggestion: String? = "Please fix your network and try again."
    }
    
    
    func fetchParkingLots() async {
        isPerformingTask = true
        defer { isPerformingTask = false }
        do {
            let metadatas = try await storekeeperClient.metadatas()
            let states = try await storekeeperClient.states()
            self.parkingLots = Dictionary(uniqueKeysWithValues: metadatas.map{ id, metadata in
                (id, ParkingLot(metadata: metadata, state: states[id]!))
            })
        } catch {
            print("error \(error)")
            self.fetchError = FetchError(
                errorDescription: error.localizedDescription
            )
            self.fetchFailed = true
        }
    }
}
