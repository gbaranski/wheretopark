package app.wheretopark.storekeeper

import app.wheretopark.shared.ParkingLotID
import app.wheretopark.shared.ParkingLotMetadata
import app.wheretopark.shared.ParkingLotState

class MemoryStore : Store {
    private val states = mutableMapOf<ParkingLotID, ParkingLotState>()
    private val metadatas = mutableMapOf<ParkingLotID, ParkingLotMetadata>()

    override suspend fun getMetadata(id: ParkingLotID) = metadatas[id]

    override suspend fun getMetadatas() = metadatas

    override suspend fun getState(id: ParkingLotID) = states[id]

    override suspend fun getStates() = states

    override suspend fun updateMetadatas(metadata: Map<ParkingLotID, ParkingLotMetadata>) {
        this.metadatas.putAll(metadata)
    }

    override suspend fun updateStates(states: Map<ParkingLotID, ParkingLotState>) {
        this.states.putAll(states)
    }
}