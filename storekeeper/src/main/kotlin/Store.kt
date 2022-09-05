package app.wheretopark.storekeeper

import app.wheretopark.shared.ParkingLotID
import app.wheretopark.shared.ParkingLotMetadata
import app.wheretopark.shared.ParkingLotState

interface Store {
    suspend fun getMetadatas(): Map<ParkingLotID, ParkingLotMetadata>
    suspend fun getStates(): Map<ParkingLotID, ParkingLotState>
    suspend fun updateMetadatas(metadata: Map<ParkingLotID, ParkingLotMetadata>)
    suspend fun updateStates(states: Map<ParkingLotID, ParkingLotState>)
}