package app.wheretopark.shared

import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.request.*
import io.ktor.http.*
import io.ktor.serialization.kotlinx.json.*
import kotlinx.serialization.json.Json

val DEFAULT_STOREKEEPER_URL = Url("https://storekeeper.wheretopark.app")

class StorekeeperClient(
    private val baseURL: Url = DEFAULT_STOREKEEPER_URL
) {
    constructor(baseURL: String = DEFAULT_STOREKEEPER_URL.toString()) : this(Url(baseURL))

    init {
        println("using `$baseURL` as Storekeeper URL")
    }

    private val http = HttpClient {
        expectSuccess = true
        install(ContentNegotiation) {
            json(Json {
                prettyPrint = true
                isLenient = true
                ignoreUnknownKeys = true
            })
        }
    }

    suspend fun metadatas(): Map<ParkingLotID, ParkingLotMetadata> {
        val response = http.get(baseURL) {
            url {
                appendPathSegments("parking-lot", "metadata")
            }
        }
        return response.body()
    }

    suspend fun postMetadatas(metadatas: Map<ParkingLotID, ParkingLotMetadata>)  {
        http.post(baseURL) {
            contentType(ContentType.Application.Json)
            setBody(metadatas)
            url {
                appendPathSegments("parking-lot", "metadata")
            }
        }
    }

    suspend fun states(): Map<ParkingLotID, ParkingLotState> {
        val response = http.get(baseURL) {
            url {
                appendPathSegments("parking-lot", "state")
            }
        }
        return response.body()
    }

    suspend fun postStates(states: Map<ParkingLotID, ParkingLotState>)  {
        http.post(baseURL) {
            contentType(ContentType.Application.Json)
            setBody(states)
            url {
                appendPathSegments("parking-lot", "state")
            }
        }
    }
}