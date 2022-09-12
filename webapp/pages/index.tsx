import {useEffect, useState} from 'react'
import Map, {MapView} from '../components/Map'
import Details from '../components/Details'
import {ParkingLotID, parseParkingLots, AuthorizationClient, StorekeeperClient} from '../lib/types'
import styles from '../styles/Home.module.css'
import List from '../components/List'
import Image from 'next/image'
import {app} from "wheretopark-shared";
import AccessType = app.wheretopark.shared.AccessType;
import {GetServerSideProps, GetServerSidePropsContext} from "next";
import getConfig from "next/config";

type HomeProps = {
    parkingLots: string
}

const Home = ({parkingLots: parkingLotsJSON}: HomeProps) => {
    const parkingLots = parseParkingLots(JSON.stringify(parkingLotsJSON))
    const [selectedParkingLotID, setSelectedParkingLotID] = useState<ParkingLotID | undefined>(undefined)
    const [mapView, setMapView] = useState<MapView>({
        width: 0,
        height: 0,
        longitude: 18.64,
        latitude: 54.35,
        zoom: 10,
        bearing: 0,
        pitch: 0,
        padding: {
            top: 0,
            bottom: 0,
            left: 0,
            right: 0,
        },
    });

    useEffect(() => {
        if (!selectedParkingLotID) return;
        const parkingLot = parkingLots[selectedParkingLotID];
        const {latitude, longitude} = parkingLot.metadata.location;
        mapView.longitude = longitude;
        mapView.latitude = latitude;
        mapView.zoom = 15;
        setMapView(mapView);
    }, [selectedParkingLotID]);

    return (
        <>
            <div className={styles.split} id={styles.master}>
                <Map parkingLots={parkingLots} mapView={mapView} selectParkingLot={setSelectedParkingLotID}/>
            </div>
            <div className={styles.split} id={styles.slave}>
                <div style={{padding: 15}}>
                    <Image alt="logo" src="/wheretopark.svg" width={100} height={14} layout={'responsive'}/>
                </div>
                <div style={{display: selectedParkingLotID ? 'none' : 'block'}}>
                    <List parkingLots={parkingLots} onSelect={setSelectedParkingLotID}/>
                </div>
                {selectedParkingLotID &&
                    <Details parkingLot={[selectedParkingLotID, parkingLots[selectedParkingLotID]]!}
                             onDismiss={() => setSelectedParkingLotID(undefined)}/>
                }
            </div>
        </>
    )
}

const {serverRuntimeConfig} = getConfig();
const authorizationClient = new AuthorizationClient(serverRuntimeConfig.AUTHORIZATION_URL, serverRuntimeConfig.CLIENT_ID, serverRuntimeConfig.CLIENT_SECRET!)
const storekeeperClient = new StorekeeperClient(serverRuntimeConfig.STOREKEEPER_URL, [AccessType.ReadMetadata, AccessType.ReadState], authorizationClient)

export async function getServerSideProps({req, res}: GetServerSidePropsContext) {
    const parkingLots = await storekeeperClient.parkingLots()
    res.setHeader(
        'Cache-Control',
        'public, s-maxage=10, stale-while-revalidate=59'
    )
    const props: HomeProps = {
        parkingLots: JSON.parse(parkingLots),
    }
    return {
        props, // will be passed to the page component as props
    }
}

export default Home
