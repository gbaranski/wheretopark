import SwiftUI

@main
struct iOSApp: App {
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
