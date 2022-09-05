package app.wheretopark.storekeeper

import app.wheretopark.shared.ParkingLotID
import app.wheretopark.shared.ParkingLotMetadata
import app.wheretopark.shared.ParkingLotState

class MemoryStore: Store {
    private val states = mutableMapOf<ParkingLotID, ParkingLotState>()
    private val metadatas = mutableMapOf<ParkingLotID, ParkingLotMetadata>()

    override suspend fun getMetadatas(): Map<ParkingLotID, ParkingLotMetadata> {
        return metadatas
    }

    override suspend fun getStates(): Map<ParkingLotID, ParkingLotState> {
        return states
    }

    override suspend fun updateMetadatas(metadata: Map<ParkingLotID, ParkingLotMetadata>) {
        this.metadatas.putAll(metadata)
    }

    override suspend fun updateStates(states: Map<ParkingLotID, ParkingLotState>) {
        this.states.putAll(states)
    }
}