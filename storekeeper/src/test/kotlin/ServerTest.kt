package app.wheretopark.storekeeper

import app.wheretopark.shared.*
import com.auth0.jwt.JWT
import com.auth0.jwt.algorithms.Algorithm
import io.ktor.client.plugins.auth.*
import io.ktor.client.plugins.auth.providers.*
import io.ktor.serialization.kotlinx.json.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.server.testing.*
import kotlinx.serialization.json.Json
import org.junit.Test
import java.net.URI
import kotlin.test.assertEquals

class ServerTest {
    private fun withClient(test: suspend (storekeeperClient: StorekeeperClient) -> Unit) = testApplication {
        val config = Config(
            jwtSecret = "someJwtSecret",
            storeURI = URI(""),
            port = 0
        )
        application {
            configure(MemoryStore(), config)
        }

        val client = createClient {
            expectSuccess = true
            install(ContentNegotiation) {
                json(Json {
                    prettyPrint = true
                    isLenient = true
                })
            }
            install(Auth) {
                bearer {
                    loadTokens {
                        val token = JWT.create().withClaim("scope", AccessType.values().toSet().encode())
                            .sign(Algorithm.HMAC512(config.jwtSecret))
                        BearerTokens(token, "")
                    }
                }
            }
        }
        val storekeeperClient = StorekeeperClient(client)
        test(storekeeperClient)
    }

    private val parkingLots = mapOf(
            (ParkingLot.galeriaBaltycka.metadata.location.hash() to ParkingLot.galeriaBaltycka),
            (ParkingLot.forumGdansk.metadata.location.hash() to ParkingLot.forumGdansk)
        )

    private val metadatas = parkingLots.split().first
    private val states = parkingLots.split().second

    @Test
    fun testParkingLotRW() = withClient {
        assertEquals(it.parkingLots().count(), 0)
        it.postParkingLots(parkingLots)
        assertEquals(it.parkingLots(), parkingLots)
        assertEquals(it.metadatas(), metadatas)
        assertEquals(it.states(), states)
        val parkingLot = parkingLots.entries.first()
        assertEquals(it.parkingLot(parkingLot.key), parkingLot.value)
        assertEquals(it.state(parkingLot.key), parkingLot.value.state)
        assertEquals(it.metadata(parkingLot.key), parkingLot.value.metadata)
    }

    @Test
    fun testStateOnlyRW() = withClient {
        assertEquals(it.states().count(), 0)
        it.postStates(states)
        assertEquals(it.states(), states)
        val parkingLot = states.entries.first()
        assertEquals(it.state(parkingLot.key), parkingLot.value)
        assertEquals(it.metadata(parkingLot.key), null)
        assertEquals(it.parkingLot(parkingLot.key), null)
    }

    @Test
    fun testMetadataOnlyRW() = withClient {
        assertEquals(it.metadatas().count(), 0)
        it.postMetadatas(metadatas)
        assertEquals(it.metadatas(), metadatas)
        val parkingLot = metadatas.entries.first()
        assertEquals(it.metadata(parkingLot.key), parkingLot.value)
        assertEquals(it.state(parkingLot.key), null)
        assertEquals(it.parkingLot(parkingLot.key), null)
    }
}
