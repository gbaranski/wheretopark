plugins {
    kotlin("jvm")
    kotlin("plugin.serialization")
}

dependencies {
    implementation(project(":shared"))

    implementation("ch.qos.logback:logback-classic:1.2.11")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.6.4")

    testImplementation(kotlin("test"))
}

