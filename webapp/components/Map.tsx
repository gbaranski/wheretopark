import ReactMapGL, {Marker} from 'react-map-gl';
import {ParkingLot, ParkingLotID} from '../lib/types';
import getConfig from "next/config";

import Image from 'next/image'

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
    selectedParkingLot: ParkingLot | null,
}

const {publicRuntimeConfig} = getConfig();

export default ({parkingLots, selectedParkingLot}: MapProps) => {
    return (
        <ReactMapGL
            mapboxAccessToken={publicRuntimeConfig.MAPBOX_ACCESS_TOKEN}
            initialViewState={{
                longitude: selectedParkingLot?.metadata.location.longitude ?? 18.64,
                latitude: selectedParkingLot?.metadata.location.latitude ?? 54.35,
                zoom: selectedParkingLot ? 15 : 10,
            }}
            style={{width: "fit"}}
            mapStyle="mapbox://styles/mapbox/navigation-day-v1"
        >
            {Object.entries(parkingLots).map(([id, parkingLot]) => (
                <MapMarker key={id} parkingLot={[id, parkingLot]}/>
            ))}
        </ReactMapGL>
    );
}