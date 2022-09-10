import { useEffect, useState } from 'react'
import type { GetServerSidePropsContext, NextPage } from 'next'
import Map, { MapView } from '../components/Map'
import Details from '../components/Details'
import { ID, ParkingLot } from '../lib/parkingLot'
import { fetchParkingLots } from '../lib/storekeeper'
import styles from '../styles/Home.module.css'
import List from '../components/List'

const Home = ({ parkingLots }: { parkingLots: ParkingLot[] }) => {
  const [selectedParkingLotID, setSelectedParkingLotID] = useState<ID | undefined>(undefined)
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
    const parkingLot = parkingLots.find((p) => p.id == selectedParkingLotID)!;
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
        <img src="wheretopark.svg" style={{ display: 'block', marginLeft: 'auto', marginRight: 'auto', width: '50%', padding: 20 }} />
        <List parkingLots={parkingLots} onSelect={setSelectedParkingLotID} />
        {selectedParkingLotID &&
          <Details parkingLot={parkingLots.find((p) => p.id == selectedParkingLotID)!} onDismiss={() => setSelectedParkingLotID(undefined)} />
        }
      </div>
    </>
  )
}

export async function getServerSideProps(context: GetServerSidePropsContext) {
  const parkingLots = await fetchParkingLots();
  return {
    props: { parkingLots }, // will be passed to the page component as props
  }
}

export default Home
