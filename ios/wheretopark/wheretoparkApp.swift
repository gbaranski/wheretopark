import SwiftUI

@main
struct iOSApp: App {
    @StateObject var appState = AppState()

	var body: some Scene {
		WindowGroup {
            ContentView()
                .environmentObject(appState)
		}
	}
}
