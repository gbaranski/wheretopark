package app.wheretopark.storekeeper

import app.wheretopark.shared.ParkingLot
import app.wheretopark.shared.ParkingLotMetadata
import app.wheretopark.shared.StorekeeperClient
import io.ktor.serialization.kotlinx.json.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.server.testing.*
import kotlinx.serialization.json.Json
import org.junit.Test
import kotlin.test.assertEquals

class ServerTest {
    private fun withClient(test: suspend (storekeeperClient: StorekeeperClient) -> Unit) = testApplication {
        application {
            configure(MemoryStore())
        }

        val client = createClient {
            install(ContentNegotiation) {
                json(Json {
                    prettyPrint = true
                    isLenient = true
                })
            }
        }
        val storekeeperClient = StorekeeperClient(client)
        test(storekeeperClient)
    }


    @Test
    fun testMetadataWrite() = withClient { it ->
        assertEquals(it.metadatas().count(), 0)
        val expected = mapOf(
            (ParkingLot.galeriaBaltycka.metadata.location.hash() to ParkingLot.galeriaBaltycka.metadata),
            (ParkingLot.forumGdansk.metadata.location.hash() to ParkingLot.forumGdansk.metadata)
        )
        it.postMetadatas(expected)
        assertEquals(it.metadatas(), expected)
    }

    @Test
    fun testStateWrite() = withClient { it ->
        assertEquals(it.states().count(), 0)
        val expected = mapOf(
            (ParkingLot.galeriaBaltycka.metadata.location.hash() to ParkingLot.galeriaBaltycka.state),
            (ParkingLot.forumGdansk.metadata.location.hash() to ParkingLot.forumGdansk.state)
        )
        it.postStates(expected)
        assertEquals(it.states(), expected)
    }
}
