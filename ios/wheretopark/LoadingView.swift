//
//  LoadingView.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 27/09/2022.
//

import SwiftUI

struct LoadingView: View {
    let description: String
    @State private var showDescription: Bool = false

    @Sendable private func delayLoadingText() async {
        try? await Task.sleep(nanoseconds: 2_000_000_000)
        showDescription = true
    }
    
    var body: some View {
        ProgressView(showDescription ? description : "")
            .task(delayLoadingText)
    }
}

struct LoadingView_Previews: PreviewProvider {
    static var previews: some View {
        LoadingView(
            description: "Loading actors"
        )
    }
}
