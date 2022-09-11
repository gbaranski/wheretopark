@file:OptIn(ExperimentalJsExport::class, DelicateCoroutinesApi::class)

package app.wheretopark.shared

import kotlinx.coroutines.DelicateCoroutinesApi
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.promise
import kotlin.js.Promise

@JsExport
class JSAuthorizationClient(url: String, clientID: String, clientSecret: String) {
    internal val client = AuthorizationClient(url, clientID, clientSecret)
}