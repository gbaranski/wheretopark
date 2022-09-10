import { Accordion, AccordionDetails, AccordionSummary, Avatar, Box, ButtonGroup, Collapse, Divider, Grid, IconButton, Link, List, ListItem, ListItemButton, ListItemIcon, ListItemText, makeStyles, rgbToHex, SxProps, Typography } from "@mui/material"
import { ParkingLot } from "../lib/parkingLot"
import { AccessTimeOutlined, ArrowBack, Call, CallOutlined, Directions, DirectionsCarOutlined, ExpandLess, ExpandMore, Favorite, Group, MailOutline, PlaceOutlined, Public, PublicOutlined, Scale, Share, Star, StarBorder, StarHalf, UpdateOutlined } from "@mui/icons-material"
import { DateTime, Duration } from 'luxon'
import { useState } from "react"
import { capitalizeFirstLetter } from "../lib/utils"
import { formatDistance, formatISO, parseISO } from "date-fns"

type Props = {
    parkingLot: ParkingLot
    onDismiss: () => void
}

const starStyle: SxProps = {
    color: "rgb(249, 176, 11)",
};

const ParkingLotDetails = ({ parkingLot, onDismiss }: Props) => {
    const openDirections = () => {
        const { latitude, longitude } = parkingLot.metadata.location;
        window.open(`https://www.google.com/maps/search/?api=1&query=${latitude}%2C${longitude}`);
    }
    const openWebsite = () => {
        const website = parkingLot.metadata.websites[0];
        if (!website) {
            // TODO: Handle it better
            return alert("No website associated");
        }
        window.open(parkingLot.metadata.websites[0]!);
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

    const [hoursOpen, setHoursOpen] = useState(false);
    const lastUpdated = formatDistance(parseISO(parkingLot.state.lastUpdated), new Date(), { addSuffix: true })

    return (
        <div>
            <Box paddingX={3} paddingY={1}>
                <Typography variant="h5" >
                    {parkingLot.metadata.name}
                </Typography>
                <Typography variant="body2">
                    <Box display="inline-block">
                        4.6
                        <span style={{ marginLeft: 5, marginRight: 5, verticalAlign: 'text-top' }}>
                            <Star fontSize="inherit" sx={starStyle} />
                            <Star fontSize="inherit" sx={starStyle} />
                            <Star fontSize="inherit" sx={starStyle} />
                            <Star fontSize="inherit" sx={starStyle} />
                            <StarHalf fontSize="inherit" sx={starStyle} />
                        </span>
                        <Link href="#" underline="hover">154 reviews</Link>
                    </Box>
                </Typography>
            </Box>
            <Divider />
            <Box display="flex" flexDirection="row" justifyContent="space-evenly" marginY={1}>
                <IconButton aria-label="show directions" onClick={openDirections}>
                    <Avatar>
                        <Directions />
                    </Avatar>
                </IconButton>
                <IconButton aria-label="call" onClick={call}>
                    <Avatar>
                        <Call />
                    </Avatar>
                </IconButton>
                <IconButton aria-label="open website" onClick={openWebsite}>
                    <Avatar>
                        <Public />
                    </Avatar>
                </IconButton>
                <IconButton aria-label="add to favourites" onClick={addToFavourites}>
                    <Avatar>
                        <Favorite />
                    </Avatar>
                </IconButton>
                <IconButton aria-label="share" onClick={share}>
                    <Avatar>
                        <Share />
                    </Avatar>
                </IconButton>
            </Box>
            <Divider />
            <List>
                <ListItem>
                    <ListItemIcon>
                        <PlaceOutlined />
                    </ListItemIcon>
                    <ListItemText>
                        {parkingLot.metadata.address}
                    </ListItemText>
                </ListItem>
                <ListItem>
                    <ListItemIcon>
                        <DirectionsCarOutlined />
                    </ListItemIcon>
                    <ListItemText primary={`${parkingLot.state.availableSpots} available spots`} />
                </ListItem>
                <ListItem>
                    <ListItemIcon>
                        <UpdateOutlined />
                    </ListItemIcon>
                    <ListItemText primary={`Last updated ${lastUpdated}`} />
                </ListItem>
                <ListItemButton onClick={() => setHoursOpen(!hoursOpen)}>
                    <ListItemIcon>
                        <AccessTimeOutlined />
                    </ListItemIcon>
                    <ListItemText >
                        <Typography display="inline" color="red">Closed</Typography>
                        <Typography display="inline"> ⋅ Opens at 12PM</Typography>
                    </ListItemText>
                    {hoursOpen ? <ExpandLess /> : <ExpandMore />}
                </ListItemButton>
                <Collapse in={hoursOpen} sx={{ pl: 9, pr: 5 }}>
                    {parkingLot.metadata.rules.map((rule) => (
                        <div>
                            <Typography variant="h6">{capitalizeFirstLetter(rule.weekdays.start.toString())} - {capitalizeFirstLetter(rule.weekdays.end.toString())}</Typography>
                            {rule.hours &&
                                <Typography variant="subtitle2">{rule.hours.start}-{rule.hours.end}</Typography>
                            }
                            {rule.pricing.map((pricing) => (
                                <div>
                                    <Box display="flex" flex-directions="row" justifyContent="space-between">
                                        <Typography display="inline" align="left">{pricing.repeating && "Each "}{Duration.fromISO(pricing.duration).toHuman()}</Typography>
                                        <Typography display="inline" align="right">{pricing.price}{parkingLot.metadata.currency}</Typography>
                                    </Box>
                                    <Divider />
                                </div>
                            ))}

                        </div>
                    ))}
                </Collapse>
                <ListItem>
                    <ListItemIcon>
                        <PublicOutlined />
                    </ListItemIcon>
                    <ListItemText >
                        {parkingLot.metadata.websites.map((url, index) => (
                            <Link href={url.toString()} underline="hover">
                                {`${url.toString()}`}{index != parkingLot.metadata.websites.length - 1 && ", "}
                            </Link>
                        ))}
                    </ListItemText>
                </ListItem>

                <ListItem>
                    <ListItemIcon>
                        <CallOutlined />
                    </ListItemIcon>
                    <ListItemText >
                        {parkingLot.metadata.phoneNumbers.map((phoneNumber, index) => (
                            <Link href={`tel:${phoneNumber}`} underline="hover">
                                {phoneNumber}{index != parkingLot.metadata.phoneNumbers.length - 1 && ", "}
                            </Link>
                        ))}
                    </ListItemText>
                </ListItem>

                <ListItem>
                    <ListItemIcon>
                        <MailOutline />
                    </ListItemIcon>
                    <ListItemText >
                        {parkingLot.metadata.emails.map((email, index) => (
                            <Link href={`mailto:${email}`} underline="hover">
                                {email}{index != parkingLot.metadata.emails.length - 1 && ", "}
                            </Link>
                        ))}
                    </ListItemText>
                </ListItem>

            </List>
            <IconButton aria-label="hide details" onClick={onDismiss}>
                <ArrowBack />
            </IconButton>
        </div>
    )
}

export default ParkingLotDetails;