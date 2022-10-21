import Details from './Details'
import {ParkingLot, ParkingLotID} from '../lib/types'
import styles from '../styles/Home.module.css'
import List from './List'
import Image from 'next/image'
import dynamic from 'next/dynamic'
import {Button, CircularProgress, Stack, Typography} from "@mui/material";
import NextLink from "next/link";
import "@fontsource/josefin-sans"
import { useState } from 'react'
import { useRouter } from 'next/router'

type HomeProps = {
    parkingLots: Record<ParkingLotID, ParkingLot>
    selectedParkingLotID?: ParkingLotID
}

const MapLoading = () => (
    <Stack style={{justifyContent: 'center', alignItems: 'center', height: '100%'}}>
        <CircularProgress/>
    </Stack>

)

const Map = dynamic(() => import("./Map"), {
    loading: MapLoading,
    ssr: false
});

export const Home = ({parkingLots, selectedParkingLotID}: HomeProps) => {
    const [selectedParkingLot, setSelectedParkingLot] = useState<[ParkingLotID, ParkingLot] | null>(selectedParkingLotID ? [selectedParkingLotID!, parkingLots[selectedParkingLotID]] : null);
    const router = useRouter()
    return (
        <>
            <div className={styles.split} id={styles.master}>
                <Button component={"a"} LinkComponent={NextLink} href="/" onClick={(e: any) => {
                    e.preventDefault();
                    setSelectedParkingLot(null);
                    router.push("/");
                }}>
                    <Typography fontSize={46} padding={1} fontFamily="Josefin Sans" fontWeight={600}  style={{textAlign: 'center'}}>
                        <span style={{color: "#313131"}}>where</span>
                        <span style={{color: "#a28a2b"}}>to</span>
                        <span style={{color: "#313131"}}>park</span>
                    </Typography>
                </Button>
                <div style={{display: selectedParkingLot ? 'none' : 'block'}}>
                    <List parkingLots={parkingLots} onSelect={(parkingLot) => {
                        setSelectedParkingLot(parkingLot);
                        router.push(`/parking-lot/${parkingLot[0]}`);
                    }}/>
                </div>
                {selectedParkingLot &&
                    <Details parkingLot={selectedParkingLot}/>
                }
            </div>
            <div className={styles.split} id={styles.slave}>
                <Map parkingLots={parkingLots} selectedParkingLot={selectedParkingLot}/>
            </div>
        </>
    )
}

export default Home
