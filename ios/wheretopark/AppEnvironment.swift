//
//  AppEnvironment.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 09/09/2022.
//

import Foundation

public struct AppEnvironment {
    enum Keys {
        static let databaseURL = "DATABASE_URL"
        static let databaseName = "DATABASE_NAME"
        static let databaseNamespace = "DATABASE_NAMESPACE"
        static let databaseUsername = "DATABASE_USERNAME"
        static let databasePassword = "DATABASE_PASSWORD"
    }
    private static let infoDictionary: [String : Any] = {
        guard let dict = Bundle.main.infoDictionary else {
            fatalError("plist file not found")
        }
        return dict
    }()
    
    static let databaseURL: URL = {
        guard let string = AppEnvironment.infoDictionary[Keys.databaseURL] as? String else {
            fatalError("database URL is not set")
        }
        return URL(string: string)!
    }()
    
    static let databaseName: String = {
        guard let string = AppEnvironment.infoDictionary[Keys.databaseName] as? String else {
            fatalError("database name is not set")
        }
        return string
    }()
    
    static let databaseNamespace: String = {
        guard let string = AppEnvironment.infoDictionary[Keys.databaseNamespace] as? String else {
            fatalError("database namespace is not set")
        }
        return string
    }()
    
    static let databaseUsername: String = {
        guard let string = AppEnvironment.infoDictionary[Keys.databaseUsername] as? String else {
            fatalError("database username is not set")
        }
        return string
    }()
    
    static let databasePassword: String = {
        guard let string = AppEnvironment.infoDictionary[Keys.databasePassword] as? String else {
            fatalError("database password is not set")
        }
        return string
    }()
}
