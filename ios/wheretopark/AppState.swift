//
//  AppState.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 11/07/2022
//

import Foundation
import SwiftUI
import Contacts

typealias ParkingLotID = String

@MainActor class AppState: ObservableObject {
    let environment = AppEnvironment()
    let databaseClient = DatabaseClient()
    
    @Published private(set) var isPerformingTask = false
    @Published var fetchError: FetchError? = nil
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
            self.parkingLots = try await databaseClient.parkingLots()
            self.fetchError = nil
            self.fetchFailed = false
        } catch {
            print("error \(error)")
            self.fetchError = FetchError(
                errorDescription: error.localizedDescription
            )
            self.fetchFailed = true
        }
    }
}
