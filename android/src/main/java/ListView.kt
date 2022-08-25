package app.wheretopark.android

import androidx.compose.foundation.ScrollState
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.*
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.input.TextFieldValue
import androidx.compose.ui.unit.dp
import app.wheretopark.shared.*

@Composable
fun ListView(parkingLotViewModel: ParkingLotViewModel) {
    val scrollState = ScrollState(initial = 0)

    Column(modifier = Modifier.verticalScroll(scrollState)) {
        parkingLotViewModel.parkingLots.forEach{ (id, parkingLot) ->
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
fun ListViewBottomSheet(parkingLotViewModel: ParkingLotViewModel, content: @Composable () -> Unit) {
    val scaffoldState = rememberBottomSheetScaffoldState()
    val searchTextState = remember { mutableStateOf(TextFieldValue("")) }

    LaunchedEffect(parkingLotViewModel.selectedParkingLotID, block = {
        if (parkingLotViewModel.selectedParkingLotID == null) {
            scaffoldState.bottomSheetState.expand()
        } else {
            scaffoldState.bottomSheetState.collapse()
        }
    })

    BottomSheetScaffold(
        sheetContent = {
            BottomSheetHandle(
                Modifier.align(Alignment.CenterHorizontally).padding(top = 10.dp)
            )
            Column(modifier = Modifier.padding(10.dp)) {
                SearchView(state = searchTextState)
                ListView(parkingLotViewModel)
            }
        },
        scaffoldState = scaffoldState,
        sheetPeekHeight = 90.dp,
    ) {
        content()
    }
}

