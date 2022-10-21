import ReactMapGL, {Marker, ViewState, ViewStateChangeEvent} from 'react-map-gl';
import {ParkingLot, ParkingLotID} from '../lib/types';
import getConfig from "next/config";

import Image from 'next/image'
import { useEffect, useState } from 'react';

type MapMarkerProps = {
    parkingLot: [ParkingLotID, ParkingLot],
}

const MapMarker = ({parkingLot}: MapMarkerProps) => (
    <Marker
        key={parkingLot[0]}
        longitude={parkingLot[1].metadata.location.longitude}
        latitude={parkingLot[1].metadata.location.latitude}
        anchor="bottom"
    >
        <Image alt="marker" src="/parking-lot-marker.png" width={48} height={48}/>
    </Marker>
)

type MapProps = {
    parkingLots: Record<ParkingLotID, ParkingLot>,
    selectedParkingLot: [ParkingLotID, ParkingLot] | null,
}

const {publicRuntimeConfig} = getConfig();

export const Map = ({parkingLots, selectedParkingLot}: MapProps) => {
    const [viewState, setViewState] = useState({
        latitude: 54.35,
        longitude: 18.64,
        zoom: 10,
    });
    useEffect(() => {
        console.log("refreshing");
        if(selectedParkingLot) {
            const {longitude, latitude} = selectedParkingLot![1].metadata.location;
            const newViewState = {
                latitude,
                longitude,
                zoom: 15,
            };
            setViewState(newViewState);
        }
    }, [selectedParkingLot]);
    return (
        <ReactMapGL
            {...viewState}
            onMove={(e: ViewStateChangeEvent) => setViewState(e.viewState)}
            mapboxAccessToken={publicRuntimeConfig.MAPBOX_ACCESS_TOKEN}
            style={{width: "fit"}}
            mapStyle="mapbox://styles/mapbox/navigation-day-v1"
        >
            {Object.entries(parkingLots).map(([id, parkingLot]) => (
                <MapMarker key={id} parkingLot={[id, parkingLot]}/>
            ))}
        </ReactMapGL>
    );
}

export default Map;