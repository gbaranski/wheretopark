package app.wheretopark.shared.client

class Greeting {
    fun greeting(): String {
        return "Hello, ${Platform().platform}!"
    }
}