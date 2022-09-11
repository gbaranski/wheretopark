import ReactMapGL, { MapProvider, Marker, useMap, ViewState } from 'react-map-gl';
import { MAPBOX_ACCESS_TOKEN } from '../environment';
import { ParkingLotID, ParkingLot } from '../lib/types';
import Image from 'next/image'

type MapMarkerProps = {
  parkingLot: [ParkingLotID, ParkingLot],
  onClick: () => void,
}

const MapMarker = ({ parkingLot, onClick }: MapMarkerProps) => (
  <Marker
    key={parkingLot[0]}
    longitude={parkingLot[1].metadata.location.longitude}
    latitude={parkingLot[1].metadata.location.latitude}
    anchor="bottom"
    onClick={onClick}
  >
    <Image alt="marker" src="/parking-lot-marker.png" width={48} height={48} />
  </Marker>
)

export type MapView = ViewState & { width: number; height: number; };

type MapProps = {
  parkingLots: Record<ParkingLotID, ParkingLot>,
  mapView: MapView,
  selectParkingLot: (id: ParkingLotID | undefined) => void
}

const Map = ({ parkingLots, selectParkingLot, mapView }: MapProps) => {
  let { mymap: map } = useMap();

  return (
    <ReactMapGL
      mapboxAccessToken={MAPBOX_ACCESS_TOKEN}
      id="mymap"
      initialViewState={{
        longitude: 18.64,
        latitude: 54.35,
        zoom: 10
      }}
      style={{ width: "fit" }}
      mapStyle="mapbox://styles/mapbox/navigation-day-v1"
    >
      {Object.entries(parkingLots).map((parkingLot) => (
        <MapMarker key={parkingLot[0]} parkingLot={parkingLot} onClick={() => selectParkingLot(parkingLot[0])} />
      ))}
    </ReactMapGL>
  );
}

const MapRoot = (props: MapProps) => (
  <MapProvider>
    <Map {...props} />
  </MapProvider>
)


export default MapRoot;