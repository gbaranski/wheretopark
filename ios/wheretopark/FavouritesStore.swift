//
//  FavouriteStore.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 10/09/2022.
//

import Foundation

class FavouritesStore {
    let store = NSUbiquitousKeyValueStore.default
    let key = "favourites"
    
    func all() -> [ParkingLotID] {
        return store.array(forKey: key)?.map { $0 as! String } ?? []
    }
    
    func exists(id: ParkingLotID) -> Bool {
        let favourites = all()
        print("favourites: \(favourites)")
        return favourites.contains(id)
    }
    
    func add(id: ParkingLotID) {
        if (exists(id: id)) { return }
        var favourites = all()
        favourites.append(id)
        print("new favourites: \(favourites)")
        store.set(favourites, forKey: key)
        store.synchronize()
        print("after update favourites: \(all())")
    }
    
    func remove(id: ParkingLotID) {
        if (!exists(id: id)) { return }
        var favourites = all()
        favourites.removeAll{ $0 == id }
        store.set(favourites, forKey: key)
    }
}
