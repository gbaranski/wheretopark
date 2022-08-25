pluginManagement {
    repositories {
        google()
        mavenCentral()
        maven(url = "https://plugins.gradle.org/m2/")
    }
}

rootProject.name = "wheretopark"

include(":shared")
include(":android")
include(":storekeeper")
include(":providers:tristar")
include(":providers:cctv")
