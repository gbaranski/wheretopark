package app.wheretopark.storekeeper

import app.wheretopark.shared.*
import io.ktor.http.*
import io.ktor.serialization.kotlinx.json.json
import io.ktor.server.application.Application
import io.ktor.server.application.call
import io.ktor.server.application.install
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.server.plugins.callloging.*
import io.ktor.server.plugins.contentnegotiation.ContentNegotiation
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.net.URI

fun main() {
    val port = System.getenv("PORT")?.toInt() ?: 8080

    embeddedServer(Netty, port = port) {
        configureRouting()
    }.start(wait = true)
}

fun Application.configureRouting() {
    install(ContentNegotiation) {
        json()
    }
    install(CallLogging)

    val storeURL = URI(System.getenv("STORE_URL") ?: "memory:")
    val store = when(storeURL.scheme) {
        "memory" -> {
            MemoryStore()
        }
        "redis" -> {
            RedisStore(storeURL.host, if (storeURL.port == -1) 6379 else storeURL.port)
        }
        else -> {
            throw IllegalArgumentException("Unknown store scheme: ${storeURL.scheme}")
        }
    }
    routing {
        get("/health-check") {
            call.respond("Hello, World!")
        }
        get("/parking-lot/state"){
            call.respond(store.getStates())
        }
        post("/parking-lot/state"){
            val newStates = call.receive<Map<ParkingLotID, ParkingLotState>>()
            store.updateStates(newStates)
            call.respondText("updated ${newStates.count()} states", status = HttpStatusCode.Accepted)
        }

        get("/parking-lot/metadata"){
            call.respond(store.getMetadatas())
        }
        post("/parking-lot/metadata"){
            val newMetadatas = call.receive<Map<ParkingLotID, ParkingLotMetadata>>()
            store.updateMetadatas(newMetadatas)
            call.respondText("updated ${newMetadatas.count()} states", status = HttpStatusCode.Accepted)
        }
    }
}