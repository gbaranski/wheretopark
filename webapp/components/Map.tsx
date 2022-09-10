import { MapsHomeWork } from '@mui/icons-material';
import { useEffect, useState } from 'react';
import ReactMapGL, { MapProvider, Marker, useMap, ViewState } from 'react-map-gl';
import { MAPBOX_ACCESS_TOKEN } from '../environment';
import { ID, ParkingLot } from '../lib/parkingLot';

type MapMarkerProps = {
  parkingLot: ParkingLot,
  onClick: () => void,
}

const MapMarker = ({ parkingLot, onClick }: MapMarkerProps) => (
  <Marker
    key={parkingLot.id}
    longitude={parkingLot.metadata.location.longitude}
    latitude={parkingLot.metadata.location.latitude}
    anchor="bottom"
    onClick={onClick}
  >
    <img src="parking-lot-marker.png" width={48} />
  </Marker>
)

export type MapView = ViewState & { width: number; height: number; };

type MapProps = {
  parkingLots: ParkingLot[],
  mapView: MapView,
  selectParkingLot: (id: ID | undefined) => void
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
      {parkingLots.map((parkingLot) => (
        <MapMarker key={parkingLot.id} parkingLot={parkingLot} onClick={() => selectParkingLot(parkingLot.id)} />
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