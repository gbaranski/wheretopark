import { Box, List, ListItem, ListItemButton, ListItemText, NativeSelect, Typography } from "@mui/material"
import { useEffect, useState } from "react";
import { ParkingLot, ID, Location } from "../lib/parkingLot"
import SearchPlace from "./SearchPlace";
import { getDistance } from 'geolib';
import { prettyDistance } from "../lib/utils";

const ParkingLotListItem = ({ parkingLot, origin, onSelect }: { parkingLot: ParkingLot, origin: Location | null, onSelect: () => void }) => (
  <ListItem key={parkingLot.id}>
    <ListItemButton onClick={onSelect}>
      <ListItemText
        primary={parkingLot.metadata.name}
        secondary={`${parkingLot.state.availableSpots} available spots ${origin ? `| ${prettyDistance(origin, parkingLot.metadata.location)} away` : ""}`}
      />
    </ListItemButton>
  </ListItem>
)

enum SortMethod {
  Name = "name",
  Distance = "distance",
}

const ParkingLotList = ({ parkingLots, onSelect }: { parkingLots: ParkingLot[], onSelect: (id: ID) => void }) => {
  const [sortMethod, setSortMethod] = useState<SortMethod>(SortMethod.Name);
  const [origin, setOrigin] = useState<Location | null>(null)

  parkingLots.sort((a, b) => {
    switch (sortMethod) {
      case SortMethod.Distance:
        if (!origin) { console.error("no origin set"); return 0; }
        const distanceA = getDistance(origin, a.metadata.location);
        const distanceB = getDistance(origin, b.metadata.location);
        return distanceA - distanceB;
      case SortMethod.Name:
        return a.id.localeCompare(b.id);
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
        {parkingLots.map((parkingLot) => ParkingLotListItem({ parkingLot, origin, onSelect: () => onSelect(parkingLot.id) }))}
      </List>
    </div>
  )
}

export default ParkingLotList;