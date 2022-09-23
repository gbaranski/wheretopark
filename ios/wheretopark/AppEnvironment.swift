//
//  Environment.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 09/09/2022.
//

import Foundation

public struct AppEnvironment {
    enum Keys {
        static let clientID = "CLIENT_ID"
        static let clientSecret = "CLIENT_SECRET"
        static let authorizationURL = "AUTHORIZATION_URL"
        static let storekeeperURL = "STOREKEEPER_URL"
    }
    private static let infoDictionary: [String : Any] = {
        guard let dict = Bundle.main.infoDictionary else {
            fatalError("plist file not found")
        }
        return dict
    }()
    static let clientID: String = {
        guard let string = AppEnvironment.infoDictionary[Keys.clientID] as? String else {
            fatalError("client id is not set")
        }
        return string
    }()
    static let clientSecret: String = {
        guard let string = AppEnvironment.infoDictionary[Keys.clientSecret] as? String else {
            fatalError("client secret is not set")
        }
        return string
    }()
    static let authorizationURL: URL = {
        guard let string = AppEnvironment.infoDictionary[Keys.authorizationURL] as? String else {
            fatalError("authorization url is not set")
        }
        return URL(string: string)!
    }()
    static let storekeeperURL: URL = {
        guard let string = AppEnvironment.infoDictionary[Keys.storekeeperURL] as? String else {
            fatalError("storekeeper url is not set")
        }
        return URL(string: string)!
    }()
}
