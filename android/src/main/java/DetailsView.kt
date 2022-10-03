package app.wheretopark.android

import android.content.Intent
import android.net.Uri
import android.text.format.DateUtils
import android.widget.Toast
import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.*
import androidx.compose.material.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Close
import androidx.compose.material.icons.filled.Favorite
import androidx.compose.material.icons.filled.MoreVert
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import app.wheretopark.shared.ParkingLot
import app.wheretopark.shared.ParkingLotMetadata
import androidx.compose.material.Icon
import java.util.*


@Composable
fun DetailsView(parkingLot: ParkingLot, onDismiss: () -> Unit) {
    val context = LocalContext.current
    var miscMenuExpanded by remember { mutableStateOf(false) }


    fun addToFavourites() {
        Toast.makeText(context, "Added to favourites", Toast.LENGTH_SHORT).show()
    }

    fun openGoogleMaps() {
        val uri =
            "http://maps.google.com/maps?daddr=${parkingLot.metadata.location.latitude},${parkingLot.metadata.location.longitude}"
        val intent = Intent(Intent.ACTION_VIEW, Uri.parse(uri))
        context.startActivity(intent)
    }

    Column {
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                parkingLot.metadata.name,
                fontWeight = FontWeight.Bold,
                style = MaterialTheme.typography.h5,
            )
            Button(onClick = onDismiss) {
                Icon(Icons.Default.Close, contentDescription = "close view")
            }
        }
        Row(modifier = Modifier.fillMaxWidth()) {
            Button(modifier = Modifier.weight(0.8f), onClick = { openGoogleMaps() }) {
                Image(
                    painterResource(id = R.drawable.ic_baseline_directions_24),
                    contentDescription = "Navigate button",
                    modifier = Modifier.size(20.dp)
                )

                Text(text = "Navigate", Modifier.padding(start = 10.dp))
            }
            Box(
                modifier = Modifier
                    .weight(0.2f)
                    .wrapContentSize(Alignment.TopEnd),
            ) {
                Button(
                    onClick = { miscMenuExpanded = true },
                    colors = ButtonDefaults.buttonColors(
                        backgroundColor = MaterialTheme.colors.secondary
                    )
                )
                {
                    Icon(Icons.Default.MoreVert, contentDescription = "show menu")
                }
                DropdownMenu(
                    expanded = miscMenuExpanded,
                    onDismissRequest = { miscMenuExpanded = false }
                ) {
                    DropdownMenuItem(onClick = { addToFavourites() }) {
                        Text(
                            "Add to favourites",
                            Modifier
                                .weight(1f)
                                .padding(end = 5.dp)
                        )
                        Icon(
                            Icons.Default.Favorite,
                            contentDescription = "add to favourites",
                            modifier = Modifier.size(24f.dp)
                        )
                    }
                }
            }
        }
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .height(IntrinsicSize.Min), horizontalArrangement = Arrangement.Center
        ) {
            Column(modifier = Modifier.weight(0.2f)) {
                Text("AVAILABILITY", fontWeight = FontWeight.ExtraBold, color = MaterialTheme.colors.secondaryVariant)
                Text("${parkingLot.state.availableSpots} cars", fontWeight = FontWeight.Bold)
            }
            Column(modifier = Modifier.weight(0.2f)) {
                Text("HOURS", fontWeight = FontWeight.ExtraBold, color = MaterialTheme.colors.secondaryVariant)
                Text("Open", fontWeight = FontWeight.Bold)
            }
            Column(modifier = Modifier.weight(0.2f)) {
                val now = Date()
                val interval =
                    DateUtils.getRelativeTimeSpanString(parkingLot.state.lastUpdated.toEpochMilliseconds(), now.time, 0)
                Text("UPDATED", fontWeight = FontWeight.ExtraBold, color = MaterialTheme.colors.secondaryVariant)
                Text("$interval", fontWeight = FontWeight.Bold)
            }
        }
        Divider(modifier = Modifier.padding(10.dp))
        Text("Pricing", style = MaterialTheme.typography.h5, fontWeight = FontWeight.Bold)
        RulesView(metadata = parkingLot.metadata)
    }
}

@Composable
fun RulesView(metadata: ParkingLotMetadata) {
    metadata.rules.map { rule ->
        rule.expandHours().map { hours ->
            Text(hours, style = MaterialTheme.typography.subtitle1, fontWeight = FontWeight.Bold)
        }
        rule.pricing.map { pricing ->
            val currency = Currency.getInstance("PLN")
            val duration = pricing.duration
            Row(horizontalArrangement = Arrangement.SpaceEvenly, modifier = Modifier.fillMaxWidth()) {
                if (pricing.repeating) {
                    Text("Each additional $duration")
                } else {
                    Text(duration.toString())
                }
                if (pricing.price == 0.0) {
                    Text("Free")
                } else {
                    Text("${pricing.price} ${currency.symbol}")
                }
            }
            Divider()

        }
    }

}

@OptIn(ExperimentalMaterialApi::class)
@Composable
fun DetailsBottomSheet(parkingLotViewModel: ParkingLotViewModel, content: @Composable () -> Unit) {
    val scaffoldState = rememberBottomSheetScaffoldState()
    var parkingLot by remember {
        mutableStateOf<ParkingLot?>(null)
    }

    LaunchedEffect(parkingLotViewModel.selectedParkingLotID) {
        if (parkingLotViewModel.selectedParkingLotID == null) {
            scaffoldState.bottomSheetState.collapse()
        } else {
            scaffoldState.bottomSheetState.expand()
            val selectedParkingLotID = parkingLotViewModel.selectedParkingLotID!!
            parkingLot = parkingLotViewModel.parkingLots[selectedParkingLotID]!!
        }
    }

    LaunchedEffect(scaffoldState.bottomSheetState.isCollapsed) {
        if (scaffoldState.bottomSheetState.isCollapsed) {
            parkingLotViewModel.selectedParkingLotID = null
        }
    }


    BottomSheetScaffold(
        sheetContent = {
            BottomSheetHandle(
                Modifier
                    .align(Alignment.CenterHorizontally)
                    .padding(top = 10.dp)
            )
            Column(modifier = Modifier.padding(10.dp)) {
                if (parkingLot != null) {
                    DetailsView(parkingLot!!, onDismiss = {
                        parkingLotViewModel.selectedParkingLotID = null
                    })
                } else {
                    CircularProgressIndicator()
                }
            }
        },
        scaffoldState = scaffoldState,
        sheetPeekHeight = 0.dp,
    ) {
        content()
    }
}

