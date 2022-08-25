package app.wheretopark.android

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.material.ExperimentalMaterialApi
import androidx.compose.runtime.*
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import app.wheretopark.android.ui.theme.WheretoparkTheme
import app.wheretopark.shared.*
import kotlinx.coroutines.launch


class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        val parkingLotViewModel = ParkingLotViewModel()
        super.onCreate(savedInstanceState)
        setContent {
            MainView(parkingLotViewModel)
        }
    }
}

@OptIn(ExperimentalMaterialApi::class)
@Composable
fun MainView(parkingLotViewModel: ParkingLotViewModel) {
    WheretoparkTheme {
        LaunchedEffect(Unit, block = {
            parkingLotViewModel.fetchParkingLots()
        })

        ParkingLotViewBottomSheet(parkingLotViewModel) {
            ListViewBottomSheet(parkingLotViewModel) {
                MapView(parkingLotViewModel)
            }
        }
    }

}

class ParkingLotViewModel: ViewModel() {
    val parkingLots = mutableStateMapOf<ParkingLotID, ParkingLot>()
    var selectedParkingLotID by mutableStateOf<ParkingLotID?>(null)
    private val storekeeperClient = StorekeeperClient()

    fun fetchParkingLots() {
        println("fetching parking lots")
        viewModelScope.launch {
            val states = storekeeperClient.states()
            val metadatas = storekeeperClient.metadatas()
            metadatas.map { (id, metadata) ->
                id to ParkingLot(metadata = metadata, state = states[id]!!)
            }.forEach{ (id, parkingLot) ->
                parkingLots[id] = parkingLot
            }
            println("retrieved ${parkingLots.count()} parking lots")
        }
    }
}

//@Preview(showBackground = true)
//@Composable
//fun DefaultPreview() {
//    WheretoparkTheme {
//        Greeting("Android")
//    }
//}