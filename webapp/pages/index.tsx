import { useEffect, useState } from 'react'
import Map, { MapView } from '../components/Map'
import Details from '../components/Details'
import { ParkingLotID, parseParkingLots, fetchParkingLots } from '../lib/types'
import styles from '../styles/Home.module.css'
import List from '../components/List'
import Image from 'next/image'

const Home = ({ parkingLots: parkingLotsJSON }: { parkingLots: any }) => {
  console.log("about to parse parking lots")
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
    const { latitude, longitude } = parkingLot.metadata.location;
    mapView.longitude = longitude;
    mapView.latitude = latitude;
    mapView.zoom = 15;
    setMapView(mapView);
  }, [selectedParkingLotID]);

  return (
    <>
      <div className={styles.split} id={styles.master}>
        <Map parkingLots={parkingLots} mapView={mapView} selectParkingLot={setSelectedParkingLotID} />
      </div>
      <div className={styles.split} id={styles.slave}>
        <div style={{padding: 15}}>
          <Image alt="logo" src="/wheretopark.svg" style={{ padding: 20 }} width={100} height={14} layout={'responsive'}/>
        </div>
        <div style={{ display: selectedParkingLotID ? 'none' : 'block' }}>
          <List parkingLots={parkingLots} onSelect={setSelectedParkingLotID} />
        </div>
        {selectedParkingLotID &&
          <Details parkingLot={[selectedParkingLotID, parkingLots[selectedParkingLotID]]!} onDismiss={() => setSelectedParkingLotID(undefined)} />
        }
      </div>
    </>
  )
}

export async function getServerSideProps() {
  const clientID = "9f75e24c-55d5-425e-8186-d8a75c4e3e85";
  const clientSecret = "b7271ed3-78ac-436f-8bbc-0222b2b26f9b";
  const parkingLots = await fetchParkingLots(clientID, clientSecret);
  return {
    props: { parkingLots: JSON.parse(parkingLots) }, // will be passed to the page component as props
  }
}

export default Home
