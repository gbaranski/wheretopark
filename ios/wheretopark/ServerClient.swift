//
//  ServerClient.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 31/08/2023.
//

import Foundation

private struct Provider: Decodable {
    let name: String;
    let url: URL;
}


@MainActor
class ServerClient: ObservableObject {
    var jsonDecoder: JSONDecoder {
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        return decoder
    }
    
    
    private func providers() async throws -> [Provider] {
        var request = URLRequest(
            url: AppEnvironment.serverURL.appending(path: "/v1/providers")
        )
        request.httpMethod = "GET"
        request.addValue("application/json", forHTTPHeaderField: "Accept")
        let (data, response) = try await URLSession.shared.data(for: request)
        let httpResponse = response as! HTTPURLResponse
        let statusCode = httpResponse.statusCode
        if statusCode != 200 {
            throw NSError(domain: "UnexpectedStatusCode", code: 1, userInfo: ["statusCode": statusCode])
        }
        let providers = try jsonDecoder.decode([Provider].self, from: data)
        return providers
    }
    
    private func parkingLotsFromProvider(provider: Provider) async throws -> [ParkingLotID : ParkingLot] {
        var request = URLRequest(
            url: provider.url.appending(path: "/parking-lots")
        )
        request.httpMethod = "GET"
        request.addValue("application/json", forHTTPHeaderField: "Accept")
        let (data, response) = try await URLSession.shared.data(for: request)
        let httpResponse = response as! HTTPURLResponse
        let statusCode = httpResponse.statusCode
        if statusCode != 200 {
            throw NSError(domain: "UnexpectedStatusCode", code: 1, userInfo: ["statusCode": statusCode])
        }
        let parkingLots = try jsonDecoder.decode([ParkingLotID : ParkingLot].self, from: data)
        return parkingLots
        
    }
    
    func parkingLots() async throws -> [ParkingLotID : ParkingLot] {
        let providers = try await self.providers()
        var parkingLots = Dictionary<ParkingLotID, ParkingLot>()
        
        await withTaskGroup(of: Void.self) { group in
            for provider in providers {
                let values = try! await self.parkingLotsFromProvider(provider: provider)
                for parkingLot in values {
                    parkingLots[parkingLot.key] = parkingLot.value
                }
            }
        }
        return parkingLots;
    }
    
}
