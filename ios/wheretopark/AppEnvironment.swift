//
//  Environment.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 09/09/2022.
//

import Foundation

public struct AppEnvironment {
    let clientID: String
    let clientSecret: String
    let storekeeperURL: URL
    let authorizationURL: URL
    
    public init() {
        let e = ProcessInfo.processInfo.environment
        clientID = e["CLIENT_ID"]!
        clientSecret = e["CLIENT_SECRET"]!
        storekeeperURL = URL(string: e["STOREKEEPER_URL"]!)!
        authorizationURL = URL(string: e["AUTHORIZATION_URL"]!)!
    }
}
