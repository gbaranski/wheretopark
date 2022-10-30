//
//  DatabaseClient.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 30/10/2022.
//

import Foundation

struct DatabaseResponse<T: Decodable>: Decodable {
    let time: String
    let status: String
    let result: [T]
}

struct DatabaseParkingLot: Decodable {
    let id: String
    let metadata: ParkingLotMetadata
    let state: ParkingLotState
}


@MainActor
class DatabaseClient: ObservableObject {
    var jsonDecoder: JSONDecoder {
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        return decoder
    }
    
    func query<T: Decodable>(string: String) async throws -> DatabaseResponse<T> {
        var request = URLRequest(
            url: AppEnvironment.databaseURL.appending(path: "/sql")
        )
        request.httpMethod = "POST"
        request.httpBody = string.data(using: String.Encoding.utf8)
        request.addValue("application/json", forHTTPHeaderField: "Accept")
        let loginString = "\(AppEnvironment.databaseUsername):\(AppEnvironment.databasePassword)"
        let loginData = loginString.data(using: String.Encoding.utf8)!
        let base64LoginString = loginData.base64EncodedString()
        request.addValue("Basic \(base64LoginString)", forHTTPHeaderField: "Authorization")
        request.addValue(AppEnvironment.databaseNamespace, forHTTPHeaderField: "NS")
        request.addValue(AppEnvironment.databaseName, forHTTPHeaderField: "DB")
        let (data, response) = try await URLSession.shared.data(for: request)
        assert((response as? HTTPURLResponse)?.statusCode == 200)
        let databaseResponse = try jsonDecoder.decode([DatabaseResponse<T>].self, from: data)
        return databaseResponse[0]
        
    }
    
    func parkingLots() async throws -> [ParkingLotID : ParkingLot] {
        let response: DatabaseResponse<DatabaseParkingLot> = try await query(string: "SELECT * FROM parking_lot")
        return Dictionary(uniqueKeysWithValues: response.result.map{ result in
            let id = result.id.components(separatedBy: ":").last!
            return (id, ParkingLot(metadata: result.metadata, state: result.state))
        })
    }
    
    
}
