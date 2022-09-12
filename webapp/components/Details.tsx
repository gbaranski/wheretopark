import {
    Avatar,
    Box,
    Collapse,
    Divider,
    IconButton,
    Link,
    List,
    ListItem,
    ListItemButton,
    ListItemIcon,
    ListItemText,
    SxProps,
    Typography
} from "@mui/material"
import {
    durationToISO,
    instantToJSDate,
    ParkingLot,
    ParkingLotID,
    ParkingLotPricingRule,
    ParkingLotResource,
    ParkingLotRule,
    toArray,
    toRecord
} from "../lib/types"
import {
    AccessTimeOutlined,
    ArrowBack,
    Call,
    Directions,
    DirectionsCarOutlined,
    ExpandLess,
    ExpandMore,
    Favorite,
    PlaceOutlined,
    Public,
    PublicOutlined,
    Share,
    Star,
    StarHalf,
    UpdateOutlined
} from "@mui/icons-material"
import {Duration} from 'luxon'
import {useState} from "react"
import {capitalizeFirstLetter} from "../lib/utils"
import {formatDistance} from "date-fns"

type Props = {
    parkingLot: [ParkingLotID, ParkingLot]
    onDismiss?: () => void
}

const starStyle: SxProps = {
    color: "rgb(249, 176, 11)",
};

const PricingRules = ({i, pricing, currency}: { i: number, currency: string, pricing: ParkingLotPricingRule[] }) => {
    return (
        <>
            {pricing.map((pricing, si) => (
                <div key={`${i}/${si}`}>
                    <Box display="flex" flex-directions="row" justifyContent="space-between">
                        <Typography display="inline"
                                    align="left">{pricing.repeating && "Each "}{Duration.fromISO(durationToISO(pricing.duration)).toHuman()}</Typography>
                        <Typography display="inline"
                                    align="right">{pricing.price}{currency}</Typography>
                    </Box>
                    <Divider/>
                </div>
            ))}
        </>

    )
}

const Rules = ({parkingLot}: { parkingLot: ParkingLot }) => {
    const [hoursOpen, setHoursOpen] = useState(false);
    return (
        <>
            <ListItemButton onClick={() => setHoursOpen(!hoursOpen)}>
                <ListItemIcon>
                    <AccessTimeOutlined/>
                </ListItemIcon>
                <ListItemText>
                    <Typography display="inline" color="red">Closed</Typography>
                    <Typography display="inline"> ⋅ Opens at 12PM</Typography>
                </ListItemText>
                {hoursOpen ? <ExpandLess/> : <ExpandMore/>}
            </ListItemButton>
            <Collapse in={hoursOpen} sx={{pl: 9, pr: 5}}>
                {toArray<ParkingLotRule>(parkingLot.metadata.rules).map((rule, i) => (
                    <div key={i}>
                        <Typography
                            variant="h6">{capitalizeFirstLetter(rule.weekdays?.start.toString() ?? "Monday")} - {capitalizeFirstLetter(rule.weekdays?.end.toString() ?? "Sunday")}</Typography>
                        {rule.hours &&
                            <Typography
                                variant="subtitle2">{rule.hours.start.toString()}-{rule.hours.end.toString()}</Typography>
                        }
                        <PricingRules i={i} pricing={toArray(rule.pricing)} currency={parkingLot.metadata.currency}/>
                    </div>
                ))}
            </Collapse>
        </>
    )
}

const Resources = ({resources}: { resources: ParkingLotResource[] }) => {
    return (
        <>
            {resources.map((resource) => (
                <ListItem key={resource.url}>
                    <ListItemIcon>
                        <PublicOutlined/>
                    </ListItemIcon>
                    <ListItemText>
                        <Link href={resource.url} underline="hover">
                            {resource.url}
                        </Link>
                    </ListItemText>
                </ListItem>
            ))}
        </>
    )
}

const AvailableSpots = ({availableSpots}: { availableSpots: Record<string, number> }) => {
    return (
        <>
            {Object.entries(availableSpots).map(([type, count]) => (
                <ListItem key={type}>
                    <ListItemIcon>
                        <DirectionsCarOutlined/>
                    </ListItemIcon>
                    <ListItemText primary={`${count} available ${type.toLowerCase()} spots`}/>
                </ListItem>
            ))}
        </>
    )

}

const Buttons = ({parkingLot}: { parkingLot: ParkingLot }) => {
    const openDirections = () => {
        const {latitude, longitude} = parkingLot.metadata.location;
        window.open(`https://www.google.com/maps/search/?api=1&query=${latitude}%2C${longitude}`);
    }
    const openWebsite = () => {
        alert("TODO!");
    }
    const addToFavourites = () => {
        alert("TODO!");
    }
    const share = () => {
        alert("TODO!");
    }
    const call = () => {
        alert("TODO!");
    }
    return (
        <Box display="flex" flexDirection="row" justifyContent="space-evenly" marginY={1}>
            <IconButton aria-label="show directions" onClick={openDirections}>
                <Avatar>
                    <Directions/>
                </Avatar>
            </IconButton>
            <IconButton aria-label="call" onClick={call}>
                <Avatar>
                    <Call/>
                </Avatar>
            </IconButton>
            <IconButton aria-label="open website" onClick={openWebsite}>
                <Avatar>
                    <Public/>
                </Avatar>
            </IconButton>
            <IconButton aria-label="add to favourites" onClick={addToFavourites}>
                <Avatar>
                    <Favorite/>
                </Avatar>
            </IconButton>
            <IconButton aria-label="share" onClick={share}>
                <Avatar>
                    <Share/>
                </Avatar>
            </IconButton>
        </Box>
    )
}

const ParkingLotDetails = ({parkingLot: [id, parkingLot], onDismiss}: Props) => {

    const lastUpdated = formatDistance(instantToJSDate(parkingLot.state.lastUpdated), new Date(), {addSuffix: true})

    return (
        <>
            <Box paddingX={3} paddingY={1}>
                <Typography variant="h5">
                    {parkingLot.metadata.name}
                </Typography>

                <Box display="inline-block">
                    <Typography variant="body2" component={'span'}>
                        4.6
                        <span style={{marginLeft: 5, marginRight: 5, verticalAlign: 'text-top'}}>
                            <Star fontSize="inherit" sx={starStyle}/>
                            <Star fontSize="inherit" sx={starStyle}/>
                            <Star fontSize="inherit" sx={starStyle}/>
                            <Star fontSize="inherit" sx={starStyle}/>
                            <StarHalf fontSize="inherit" sx={starStyle}/>
                        </span>
                        <Link href="#" underline="hover">154 reviews</Link>
                    </Typography>
                </Box>
            </Box>
            <Divider/>
            <Buttons parkingLot={parkingLot}/>
            <Divider/>
            <List>
                <ListItem>
                    <ListItemIcon>
                        <PlaceOutlined/>
                    </ListItemIcon>
                    <ListItemText>
                        {parkingLot.metadata.address}
                    </ListItemText>
                </ListItem>
                <AvailableSpots availableSpots={toRecord(parkingLot.state.availableSpots)}/>
                <Rules parkingLot={parkingLot}/>
                <ListItem>
                    <ListItemIcon>
                        <UpdateOutlined/>
                    </ListItemIcon>
                    <ListItemText primary={`Last updated ${lastUpdated}`}/>
                </ListItem>
                <Resources resources={toArray(parkingLot.metadata.resources)}/>
            </List>
            {onDismiss &&
                <IconButton aria-label="hide details" onClick={onDismiss}>
                    <ArrowBack/>
                </IconButton>
            }
        </>
    )
}

export default ParkingLotDetails;