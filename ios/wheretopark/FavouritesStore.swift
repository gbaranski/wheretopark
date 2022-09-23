//
//  FavouriteStore.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 10/09/2022.
//

import Foundation
import SwiftUI

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
        store.synchronize()
    }
}


class FavouriteManager: ObservableObject {
    static let store = FavouritesStore()
    @Binding var id: ParkingLotID?
    
    @Published var isFavourite: Bool
    
    init(id: Binding<ParkingLotID?>) {
        self._id = id
        self.isFavourite = (id.wrappedValue != nil) ? Self.store.exists(id: id.wrappedValue!) : false
    }
    
    func add() {
        Self.store.add(id: id!)
        self.isFavourite = true
        self.objectWillChange.send()
    }
    
    func remove() {
        Self.store.remove(id: id!)
        self.isFavourite = false
        self.objectWillChange.send()
    }
}
