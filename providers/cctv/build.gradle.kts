
plugins {
    kotlin("jvm")
    kotlin("plugin.serialization")
    id("io.ktor.plugin") version "2.1.0"
    application
}

val ktorVersion: String by project

application {
    mainClass.set("app.wheretopark.providers.cctv.ApplicationKt")
}

dependencies {
    implementation(project(":shared"))
    implementation(project(":providers:shared"))

    implementation("ch.qos.logback:logback-classic:1.2.11")
    implementation("io.ktor:ktor-client-core:$ktorVersion")
    implementation("io.ktor:ktor-client-cio:$ktorVersion")
    implementation("io.ktor:ktor-client-content-negotiation:$ktorVersion")
    implementation("io.ktor:ktor-serialization-kotlinx-json:$ktorVersion")
    implementation("org.jetbrains.kotlinx:kotlinx-datetime:0.4.0")

    testImplementation(kotlin("test"))
}

