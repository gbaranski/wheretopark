package app.wheretopark.android

import androidx.compose.foundation.layout.*
import androidx.compose.material.Button
import androidx.compose.material.FloatingActionButton
import androidx.compose.material.Icon
import androidx.compose.material.Text
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Refresh
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import androidx.compose.ui.zIndex
import app.wheretopark.shared.ParkingLot
import app.wheretopark.shared.ParkingLotID
import com.google.accompanist.permissions.ExperimentalPermissionsApi
import com.google.android.gms.maps.CameraUpdateFactory
import com.google.android.gms.maps.model.BitmapDescriptorFactory
import com.google.android.gms.maps.model.CameraPosition
import com.google.android.gms.maps.model.LatLng
import com.google.android.gms.maps.model.MapStyleOptions
import com.google.maps.android.compose.*
import com.google.maps.android.ui.IconGenerator
import com.google.accompanist.permissions.rememberMultiplePermissionsState


@OptIn(ExperimentalPermissionsApi::class)
@Composable
fun MapView(parkingLotViewModel: ParkingLotViewModel) {
    val gdansk = LatLng(54.3520, 18.6466)
    val context = LocalContext.current
    val iconGenerator = IconGenerator(context)
    val cameraPositionState = rememberCameraPositionState {
        position = CameraPosition.fromLatLngZoom(gdansk, 10f)
    }
    val uiSettings by remember { mutableStateOf(MapUiSettings(myLocationButtonEnabled = true, zoomControlsEnabled = false)) }
    val properties by remember {
        mutableStateOf(
            MapProperties(
                mapStyleOptions = MapStyleOptions.loadRawResourceStyle(context, R.raw.map_style),
                isMyLocationEnabled = true
            )
        )
    }
    val locationPermissionsState = rememberMultiplePermissionsState(
        listOf(
            android.Manifest.permission.ACCESS_COARSE_LOCATION,
            android.Manifest.permission.ACCESS_FINE_LOCATION,
        )
    )
    if (locationPermissionsState.allPermissionsGranted) {
        LaunchedEffect(parkingLotViewModel.selectedParkingLotID, block = {
            val parkingLot =
                parkingLotViewModel.parkingLots[parkingLotViewModel.selectedParkingLotID]
                    ?: return@LaunchedEffect
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

        Box {
            Column(
                modifier = Modifier
                    .zIndex(2f)
                    .align(Alignment.BottomEnd)
                    .padding(12.dp),
                horizontalAlignment = Alignment.End,
                verticalArrangement = Arrangement.Bottom,
            ) {
                FloatingActionButton(onClick = {
                    parkingLotViewModel.fetchParkingLots()
                }) {
                    Icon(Icons.Default.Refresh, contentDescription = "refresh parking lots")
                }
//                Spacer(modifier = Modifier.height(8.dp))
//                FloatingActionButton(onClick = {}) {
//                    Icon(Icons.Default.LocationOn, contentDescription = "refresh parking lots")
//                }
            }
            GoogleMap(
                modifier = Modifier
                    .zIndex(1f),
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
    } else {
        Column(
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center,
            modifier = Modifier.fillMaxSize()
        ) {
            Text("No location permission granted")
            Button(onClick = {
                locationPermissionsState.launchMultiplePermissionRequest()
            }) {
                Text("Request permission")
            }
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