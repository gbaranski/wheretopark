import Details from './Details'
import {ParkingLot, ParkingLotID} from '../lib/types'
import styles from '../styles/Home.module.css'
import List from './List'
import Image from 'next/image'
import dynamic from 'next/dynamic'
import {CircularProgress, Stack} from "@mui/material";
import Link from "next/link";

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
    const selectedParkingLot = selectedParkingLotID ? parkingLots[selectedParkingLotID] : null;
    return (
        <>
            <div className={styles.split} id={styles.master}>
                <div style={{padding: 15}}>
                    <Link href={`/`}>
                        <a>
                            <Image alt="logo" src="/wheretopark.svg" width={100} height={14} layout={'responsive'}
                                   objectFit="contain"/>
                        </a>
                    </Link>
                </div>
                <div style={{display: selectedParkingLotID ? 'none' : 'block'}}>
                    <List parkingLots={parkingLots}/>
                </div>
                {selectedParkingLotID &&
                    <Details parkingLot={[selectedParkingLotID, parkingLots[selectedParkingLotID]]!}/>
                }
            </div>
            <div className={styles.split} id={styles.slave}>
                <Map parkingLots={parkingLots} selectedParkingLot={selectedParkingLot}/>
            </div>
        </>
    )
}

export default Home
