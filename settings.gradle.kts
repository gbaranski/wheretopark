pluginManagement {
    repositories {
        google()
        mavenCentral()
        maven(url = "https://plugins.gradle.org/m2/")
    }
}

rootProject.name = "wheretopark"

include(":shared")
include(":storekeeper")
include(":providers:tristar")
include(":providers:cctv")
include(":providers:shared")
include(":android")
