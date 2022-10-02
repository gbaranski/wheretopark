package app.wheretopark.android

import androidx.compose.foundation.layout.fillMaxHeight
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import app.wheretopark.shared.ParkingLot
import app.wheretopark.shared.ParkingLotID
import com.google.android.gms.maps.CameraUpdate
import com.google.android.gms.maps.CameraUpdateFactory
import com.google.android.gms.maps.model.BitmapDescriptorFactory
import com.google.android.gms.maps.model.CameraPosition
import com.google.android.gms.maps.model.LatLng
import com.google.android.gms.maps.model.MapStyleOptions
import com.google.maps.android.compose.*
import com.google.maps.android.ui.IconGenerator

@Composable
fun MapView(parkingLotViewModel: ParkingLotViewModel) {
    val gdansk = LatLng(54.3520, 18.6466)
    val context = LocalContext.current
    val iconGenerator = IconGenerator(context)
    val cameraPositionState = rememberCameraPositionState {
        position = CameraPosition.fromLatLngZoom(gdansk, 10f)
    }
    val uiSettings by remember { mutableStateOf(MapUiSettings(zoomControlsEnabled = false)) }
    val properties by remember {
        mutableStateOf(
            MapProperties(
                mapStyleOptions = MapStyleOptions.loadRawResourceStyle(context, R.raw.map_style)
            )
        )
    }

    LaunchedEffect(parkingLotViewModel.selectedParkingLotID, block = {
        val parkingLot =
            parkingLotViewModel.parkingLots[parkingLotViewModel.selectedParkingLotID] ?: return@LaunchedEffect
        val cameraPosition = CameraPosition.fromLatLngZoom(
            LatLng(
                parkingLot.metadata.location.latitude,
                parkingLot.metadata.location.longitude,
            ),
            15f,
        )
        val cameraUpdate = CameraUpdateFactory.newCameraPosition(cameraPosition)
        cameraPositionState.animate(cameraUpdate, 1000)
    })

    GoogleMap(
        modifier = Modifier.fillMaxWidth().fillMaxHeight(0.6F),
        cameraPositionState = cameraPositionState,
        uiSettings = uiSettings,
        properties = properties
    ) {
        parkingLotViewModel.parkingLots.map { (id, parkingLot) ->
            println("rendering parking lot: $id")
            MapMarkerView(iconGenerator, parkingLot, parkingLotID = id, parkingLotViewModel)
        }
    }
}

@Composable
fun MapMarkerView(
    iconGenerator: IconGenerator,
    parkingLot: ParkingLot,
    parkingLotID: ParkingLotID,
    parkingLotViewModel: ParkingLotViewModel
) {
    val position = LatLng(
        parkingLot.metadata.location.latitude,
        parkingLot.metadata.location.longitude
    )

    val state = rememberMarkerState(position = position)

    LaunchedEffect(parkingLotViewModel.selectedParkingLotID, block = {
        println("updated selected parking lot")
        if (parkingLotID == parkingLotViewModel.selectedParkingLotID) {
            state.showInfoWindow()
        } else {
            state.hideInfoWindow()
        }
    })

    Marker(
        icon = BitmapDescriptorFactory.fromBitmap(iconGenerator.makeIcon("P")),
        state = state,
        title = parkingLot.metadata.name,
        snippet = "${parkingLot.state.availableSpots} available parking spots",
        onClick = { _ ->
            parkingLotViewModel.selectedParkingLotID = parkingLotID
            true
        }
    )
}