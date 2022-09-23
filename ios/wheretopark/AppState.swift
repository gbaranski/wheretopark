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
    let authorizationClient: AuthorizationClient
    let storekeeperClient: StorekeeperClient
    init() {
        self.authorizationClient = AuthorizationClient(baseURL: AppEnvironment.authorizationURL.absoluteString, clientID: AppEnvironment.clientID, clientSecret: AppEnvironment.clientSecret)
        self.storekeeperClient = StorekeeperClient(
            baseURL: AppEnvironment.storekeeperURL.absoluteString,
            authorizationClient: authorizationClient,
            accessScope: Set([AccessType.readmetadata, AccessType.readstate])
        )
    }
    @Published private(set) var isPerformingTask = false
    @Published private(set) var fetchError: FetchError? = nil
    @Published var fetchFailed = false
    
    @Published private(set) var parkingLots = [ParkingLotID : ParkingLot]()
    @Published var selectedParkingLotID: ParkingLotID? = nil
    var isSelected: Binding<Bool> {
        Binding {
            self.selectedParkingLotID != nil
        } set: { value in
            if !value {
                self.selectedParkingLotID = nil
            } else {
                fatalError("unexpected value \(String(describing: value))")
            }
        }
    }
    var selectedParkingLot: Binding<ParkingLot?> {
        Binding {
            self.selectedParkingLotID != nil ? self.parkingLots[self.selectedParkingLotID!] : nil
        } set: { value in
            if value == nil {
                self.selectedParkingLotID = nil
            } else {
                fatalError("unexpected value \(String(describing: value))")
            }
        }
        
    }
    
    
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
            self.parkingLots = Dictionary(uniqueKeysWithValues: metadatas.compactMap{ id, metadata in
                let state = states[id]
                if state == nil {
                    return nil
                } else {
                    return (id, ParkingLot(metadata: metadata, state: state!))
                }
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
