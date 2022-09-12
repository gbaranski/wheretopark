package app.wheretopark.providers.tristar.gdynia

import app.wheretopark.providers.shared.Provider
import app.wheretopark.shared.*
import com.charleskorn.kaml.Yaml
import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.request.*
import io.ktor.serialization.kotlinx.json.*
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.decodeFromString
import kotlinx.serialization.json.Json
import kotlin.time.Duration
import kotlin.time.Duration.Companion.days
import kotlin.time.Duration.Companion.seconds

const val METADATA_URL = "http://api.zdiz.gdynia.pl/ri/rest/parkings"
const val STATE_URL = "http://api.zdiz.gdynia.pl/ri/rest/parking_places"

@Serializable
data class ParkingLotConfiguration(
    val resources: List<ParkingLotResource>,
    @SerialName("total-spots")
    val totalSpots: UInt,
    val features: List<ParkingLotFeature>,
    val rules: List<ParkingLotRule>
)

@Serializable
data class Configuration(
    @SerialName("parking-lots")
    val parkingLots: Map<UInt, ParkingLotConfiguration>,
)

class TristarGdyniaProvider : Provider() {
    override val name: String
        get() = "tristar/gdynia"
    override val metadataInterval: Duration
        get() = 1.days
    override val stateInterval: Duration
        get() = 30.seconds

    private val client = HttpClient {
        install(ContentNegotiation) {
            json(Json {
                prettyPrint = true
                isLenient = true
                ignoreUnknownKeys = true
            })
        }
    }

    // Mapping of vendor and storekeeper ID's
    private val ids = mutableMapOf<UInt, ParkingLotID>()

    private val configuration =
        Yaml().decodeFromString<Configuration>(this.javaClass.getResource("/configuration-gdynia.yaml")!!.readText())

    override suspend fun metadatas(): Map<ParkingLotID, ParkingLotMetadata> {
        val vendorMetadatas = client.get(METADATA_URL).body<MetadataResponse>()
        return vendorMetadatas.parkingLots.map {
            val configuration = configuration.parkingLots[it.id] ?: return@map null
            val id = it.location.coordinates.hash()
            ids[it.id] = id
            val metadata = ParkingLotMetadata(
                name = it.name,
                address = it.address,
                location = it.location.coordinates,
                resources = configuration.resources,
                totalSpots = mapOf(
                    Pair(ParkingSpotType.CAR, configuration.totalSpots)
                ),
                features = configuration.features,
                currency = "PLN",
                rules = configuration.rules,
            )
            id to metadata
        }
            .filterNotNull()
            .toMap()
    }

    override suspend fun states(): Map<ParkingLotID, ParkingLotState> {
        val vendorStates = client.get(STATE_URL).body<StateResponse>()
        return vendorStates.map {
            val id = ids[it.parkingID] ?: return@map null
            val state = ParkingLotState(
                lastUpdated = it.insertTime,
                availableSpots = mapOf(
                    Pair(ParkingSpotType.CAR, it.freePlaces)
                ),
            )
            id to state
        }
            .filterNotNull()
            .toMap()
    }
}
