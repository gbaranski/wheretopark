plugins {
    kotlin("jvm")
    kotlin("plugin.serialization")
}

val ktorVersion: String by project

dependencies {
    implementation(project(":shared"))

    implementation("io.ktor:ktor-client-core:$ktorVersion")
    implementation("ch.qos.logback:logback-classic:1.2.11")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.6.4")

    testImplementation(kotlin("test"))
}

