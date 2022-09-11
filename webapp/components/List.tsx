import { Box, List, ListItem, ListItemButton, ListItemText, NativeSelect, Typography } from "@mui/material"
import { useState } from "react";
import { ParkingLot, ParkingLotID, Coordinate } from "../lib/types"
import SearchPlace from "./SearchPlace";
import { getDistance } from 'geolib';
import { prettyDistance } from "../lib/utils";

const ParkingLotListItem = ({ parkingLot, origin, onSelect }: { parkingLot: [ParkingLotID, ParkingLot], origin: Coordinate | null, onSelect: () => void }) => (
  <ListItem key={parkingLot[0]}>
    <ListItemButton onClick={onSelect}>
      <ListItemText
        primary={parkingLot[1].metadata.name}
        secondary={`${parkingLot[1].state.availableSpots} available spots ${origin ? `| ${prettyDistance(origin, parkingLot[1].metadata.location)} away` : ""}`}
      />
    </ListItemButton>
  </ListItem>
)

enum SortMethod {
  Name = "name",
  Distance = "distance",
}

const ParkingLotList = ({ parkingLots, onSelect }: { parkingLots: Record<ParkingLotID, ParkingLot>, onSelect: (id: ParkingLotID) => void }) => {
  const [sortMethod, setSortMethod] = useState<SortMethod>(SortMethod.Name);
  const [origin, setOrigin] = useState<Coordinate | null>(null)

  Object.entries(parkingLots).sort((a, b) => {
    switch (sortMethod) {
      case SortMethod.Distance:
        if (!origin) { console.error("no origin set"); return 0; }
        const distanceA = getDistance(origin, a[1].metadata.location);
        const distanceB = getDistance(origin, b[1].metadata.location);
        return distanceA - distanceB;
      case SortMethod.Name:
        return a[0].localeCompare(b[0]);
    }

  });
  return (
    <div>
      <SearchPlace onSelect={(location) => {
        console.log("selected: ", location)
        setOrigin(location);
        setSortMethod(location ? SortMethod.Distance : SortMethod.Name);
      }} buttonNeighbour={() =>
        <Box>
          <Typography display="inline">Sort by: </Typography>

          <NativeSelect
            value={sortMethod}
            onChange={(event) => setSortMethod(event.target.value as SortMethod)}
          >
            <option value={SortMethod.Name as string}>Name</option>
            <option value={SortMethod.Distance as string} disabled={origin == null}>Distance</option>
          </NativeSelect>

        </Box>
      } />
      <List>
        {Object.entries(parkingLots).map((parkingLot) => ParkingLotListItem({ parkingLot, origin, onSelect: () => onSelect(parkingLot[0]) }))}
      </List>
    </div>
  )
}

export default ParkingLotList;