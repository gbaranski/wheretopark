package app.wheretopark.android

import androidx.compose.foundation.ScrollState
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Favorite
import androidx.compose.material.icons.filled.Refresh
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalConfiguration
import androidx.compose.ui.text.input.TextFieldValue
import androidx.compose.ui.unit.dp
import app.wheretopark.shared.ParkingLot
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.launch

@Composable
fun ListView(parkingLotViewModel: ParkingLotViewModel) {
    val scrollState = ScrollState(initial = 0)

    Column(modifier = Modifier.verticalScroll(scrollState)) {
        parkingLotViewModel.parkingLots.forEach { (id, parkingLot) ->
            RowView(
                parkingLot = parkingLot,
                onClick = {
                    parkingLotViewModel.selectedParkingLotID = id
                }
            )
        }
    }
}

@Composable
fun RowView(parkingLot: ParkingLot, onClick: () -> Unit) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 10.dp)
            .clickable { onClick() },
        elevation = 10.dp
    ) {
        Column(modifier = Modifier.padding(15.dp)) {
            Text(parkingLot.metadata.name)
            Text("${parkingLot.state.availableSpots} available parking spots")
        }
    }
}

@OptIn(ExperimentalMaterialApi::class)
@Composable
fun ListBottomSheet(parkingLotViewModel: ParkingLotViewModel, content: @Composable () -> Unit) {
    val scaffoldState = rememberBottomSheetScaffoldState()
    val searchTextState = remember { mutableStateOf(TextFieldValue("")) }
    val coroutineScope = rememberCoroutineScope()
    val configuration = LocalConfiguration.current

    LaunchedEffect(parkingLotViewModel.selectedParkingLotID, block = {
        if (parkingLotViewModel.selectedParkingLotID != null) {
            scaffoldState.bottomSheetState.collapse()
        }
    })

    BottomSheetScaffold(
        sheetContent = {
            BottomSheetHandle(
                Modifier
                    .align(Alignment.CenterHorizontally)
                    .padding(top = 10.dp)
            )
            Column(modifier = Modifier.padding(10.dp)) {
                SearchView(state = searchTextState, onTextFieldFocus = { state ->
                    if (state.hasFocus) {
                        coroutineScope.launch {
                            scaffoldState.bottomSheetState.expand()
                        }
                    }
                })
                if (parkingLotViewModel.isProcessing.value) {
                    Box(
                        contentAlignment = Alignment.Center,
                        modifier = Modifier.fillMaxWidth().padding(top = 36.dp)
                    ) {
                        CircularProgressIndicator()
                    }
                } else {
                    ListView(parkingLotViewModel)
                }
            }
        },
        scaffoldState = scaffoldState,
        sheetPeekHeight = (configuration.screenHeightDp * 0.4).dp,
    ) {
        content()
    }
}

