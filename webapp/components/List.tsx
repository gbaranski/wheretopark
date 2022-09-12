import {Box, List, ListItem, ListItemButton, ListItemText, NativeSelect, Typography} from "@mui/material"
import {useState} from "react";
import {ParkingLot, ParkingLotID, Coordinate} from "../lib/types"
import SearchPlace from "./SearchPlace";
import {prettyDistance} from "../lib/utils";
import Link from "next/link";

type ListItemProps = { parkingLot: [ParkingLotID, ParkingLot], origin: Coordinate | null }

const ParkingLotListItem = ({parkingLot: [id, parkingLot], origin}: ListItemProps) => (
    <ListItem key={id}>
        <Link href={`/parking-lot/${id}`} passHref={true}>
            <ListItemButton>
                <ListItemText
                    primary={parkingLot.metadata.name}
                    secondary={`${parkingLot.state.availableSpots} available spots ${origin ? `| ${prettyDistance(origin, parkingLot.metadata.location)} away` : ""}`}
                />
            </ListItemButton>
        </Link>
    </ListItem>
)

enum SortMethod {
    Name = "name",
    Distance = "distance",
}

type Props = {
    parkingLots: Record<ParkingLotID, ParkingLot>
}

const ParkingLotList = ({parkingLots}: Props) => {
    const [sortMethod, setSortMethod] = useState<SortMethod>(SortMethod.Name);
    const [origin, setOrigin] = useState<Coordinate | null>(null)
    const parkingLotsEntries = Object.entries(parkingLots)
    parkingLotsEntries.sort((a, b) => {
        switch (sortMethod) {
            case SortMethod.Distance:
                if (!origin) {
                    console.error("no origin set");
                    return 0;
                }
                const distanceA = origin.distanceTo(a[1].metadata.location)
                const distanceB = origin.distanceTo(b[1].metadata.location)
                console.log({distanceA, distanceB})
                return distanceA - distanceB;
            case SortMethod.Name:
                return a[0].localeCompare(b[0]);
        }
    });
    parkingLots = Object.fromEntries(parkingLotsEntries)
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
            }/>
            <List>
                {Object.entries(parkingLots).map((parkingLot) => ParkingLotListItem({
                    parkingLot,
                    origin,
                }))}
            </List>
        </div>
    )
}

export default ParkingLotList;