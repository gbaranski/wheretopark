import Details from '../components/Details'
import {ParkingLotID, ParkingLot} from '../lib/types'
import styles from '../styles/Home.module.css'
import List from '../components/List'
import Image from 'next/image'
import dynamic from 'next/dynamic'
import {CircularProgress, Stack} from "@mui/material";

type HomeProps = {
    parkingLots: Record<ParkingLotID, ParkingLot>
    selectedParkingLotID?: ParkingLotID
}

const MapLoading = () => (
    <Stack style={{justifyContent: 'center', alignItems: 'center', height: '100%'}}>
        <CircularProgress/>
    </Stack>

)

const Map = dynamic(() => import("../components/Map"), {
    loading: MapLoading,
    ssr: false
});

export const Home = ({parkingLots, selectedParkingLotID}: HomeProps) => {
    const selectedParkingLot = selectedParkingLotID ? parkingLots[selectedParkingLotID] : null;
    return (
        <>
            <div className={styles.split} id={styles.master}>
                <Map parkingLots={parkingLots} selectedParkingLot={selectedParkingLot}/>
            </div>
            <div className={styles.split} id={styles.slave}>
                <div style={{padding: 15}}>
                    <Image alt="logo" src="/wheretopark.svg" width={100} height={14} layout={'responsive'}/>
                </div>
                <div style={{display: selectedParkingLotID ? 'none' : 'block'}}>
                    <List parkingLots={parkingLots}/>
                </div>
                {selectedParkingLotID &&
                    <Details parkingLot={[selectedParkingLotID, parkingLots[selectedParkingLotID]]!}/>
                }
            </div>
        </>
    )
}

