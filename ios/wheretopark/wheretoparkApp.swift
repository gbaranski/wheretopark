//
//  wheretoparkApp.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 17/05/2022.
//

import SwiftUI

@main
struct wheretoparkApp: App {
    @StateObject var appState = AppState()
    @StateObject var locationManager = LocationManager()
    
    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(appState)
                .environmentObject(locationManager)
        }
    }
}
