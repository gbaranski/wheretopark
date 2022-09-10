package app.wheretopark.shared

import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.plugins.*
import io.ktor.client.plugins.auth.*
import io.ktor.client.plugins.auth.providers.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.request.*
import io.ktor.http.*
import io.ktor.serialization.kotlinx.json.*
import kotlinx.serialization.json.Json

const val DEFAULT_STOREKEEPER_URL = "https://storekeeper.wheretopark.app"

class StorekeeperClient(
    private val http: HttpClient,
) {
    constructor(
        baseURL: String = DEFAULT_STOREKEEPER_URL,
        authorizationClient: AuthorizationClient,
        accessScope: Set<AccessType>
    ) : this(
        HttpClient {
            defaultRequest {
                url(baseURL)
            }
            expectSuccess = true
            install(ContentNegotiation) {
                json(Json {
                    prettyPrint = true
                    isLenient = true
                    ignoreUnknownKeys = true
                })
            }
            install(HttpRequestRetry) {
                retryOnServerErrors(maxRetries = 5)
                exponentialDelay()
            }
            install(Auth) {
                bearer {
                    loadTokens {
                        val tokenResponse = authorizationClient.token(accessScope)
                        println("obtained new token: ${tokenResponse.accessToken}")
                        BearerTokens(tokenResponse.accessToken, "")
                    }
                    refreshTokens {
                        val tokenResponse = authorizationClient.token(accessScope)
                        println("obtained new token: ${tokenResponse.accessToken}")
                        BearerTokens(tokenResponse.accessToken, "")
                    }
                }
            }
        },
    )

    suspend fun metadatas(): Map<ParkingLotID, ParkingLotMetadata> = http.get("/parking-lot/metadata").body()

    suspend fun postMetadatas(metadatas: Map<ParkingLotID, ParkingLotMetadata>) =
        http.post("/parking-lot/metadata") {
            contentType(ContentType.Application.Json)
            setBody(metadatas)
        }

    suspend fun states(): Map<ParkingLotID, ParkingLotState> = http.get("/parking-lot/state").body()

    suspend fun postStates(states: Map<ParkingLotID, ParkingLotState>) =
        http.post("/parking-lot/state") {
            contentType(ContentType.Application.Json)
            setBody(states)
        }
}