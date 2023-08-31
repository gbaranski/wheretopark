//
//  AppEnvironment.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 09/09/2022.
//

import Foundation

public struct AppEnvironment {
    enum Keys {
        static let serverURL = "SERVER_URL"
    }
    private static let infoDictionary: [String : Any] = {
        guard let dict = Bundle.main.infoDictionary else {
            fatalError("plist file not found")
        }
        return dict
    }()
    
    static let serverURL: URL = {
        guard let string = AppEnvironment.infoDictionary[Keys.serverURL] as? String else {
            fatalError("server URL is not set")
        }
        return URL(string: string)!
    }()
}
