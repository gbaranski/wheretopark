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

    val states = mutableMapOf<ParkingLotID, ParkingLotState>()
    val metadatas = mutableMapOf<ParkingLotID, ParkingLotMetadata>()
    routing {
        get("/health-check") {
            call.respond("Hello, World!")
        }
        get("/parking-lot/state"){
            call.respond(states)
        }
        post("/parking-lot/state"){
            val newStates = call.receive<Map<ParkingLotID, ParkingLotState>>()
            states.putAll(newStates)
            call.respondText("updated ${newStates.count()} states", status = HttpStatusCode.Accepted)
        }

        get("/parking-lot/metadata"){
            call.respond(metadatas)
        }
        post("/parking-lot/metadata"){
            val newMetadatas = call.receive<Map<ParkingLotID, ParkingLotMetadata>>()
            metadatas.putAll(newMetadatas)
            call.respondText("updated ${newMetadatas.count()} states", status = HttpStatusCode.Accepted)
        }
    }
}