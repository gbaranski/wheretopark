
plugins {
    kotlin("jvm")
    kotlin("plugin.serialization")
    id("io.ktor.plugin") version "2.1.0"
    application
}

val ktorVersion = "2.1.0"

application {
    mainClass.set("app.wheretopark.storekeeper.ApplicationKt")
}

dependencies {
    implementation(project(":shared"))

    implementation("io.ktor:ktor-server-netty:$ktorVersion")
    implementation("io.ktor:ktor-server-content-negotiation:$ktorVersion")
    implementation("io.ktor:ktor-serialization-kotlinx-json:$ktorVersion")
    implementation("io.ktor:ktor-server-call-logging:$ktorVersion")

    implementation("ch.qos.logback:logback-classic:1.2.11")

    testImplementation(kotlin("test"))
    testImplementation("io.ktor:ktor-server-tests:$ktorVersion")
}