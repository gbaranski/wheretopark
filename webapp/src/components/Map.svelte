<script lang="ts">
    import { PUBLIC_MAPBOX_ACCESS_TOKEN } from '$env/static/public';
	import { onMount } from "svelte";
    import 'mapbox-gl/dist/mapbox-gl.css';
	import { SpotType, type ID, type ParkingLot } from "$lib/types";
	import { currentMap, parkingLots as parkingLotsStore } from "$lib/store";
	import { markerColor, parkingLotStatus } from "$lib/utils";
    import mapboxgl from "mapbox-gl";
    
    let map: mapboxgl.Map;
    $: parkingLots = $parkingLotsStore;
    $: {
        Object.entries(parkingLots).map(([id, parkingLot]) => {
            const [longitude, latitude] = parkingLot.metadata.geometry.coordinates;
            const popupHtml = `
                <h4>${parkingLot.metadata.name}</h4>
                <p>${parkingLot.state.availableSpots["CAR"]} available car spots</p>
                <a href="/parking-lot/${id}">Open</a>
            `;
            const popup = new mapboxgl.Popup({offset: 25}).setHTML(popupHtml);
            const status = parkingLotStatus(parkingLot, SpotType.CAR)[0];
            const options = {
                color: markerColor(parkingLot.state.availableSpots["CAR"], parkingLot.metadata.totalSpots["CAR"], status)
            }
            return new mapboxgl.Marker(options)
                .setLngLat([longitude, latitude])
                .setPopup(popup);
        }).forEach((marker) => marker.addTo(map));
    }

    onMount(async () => {
        map = new mapboxgl.Map({
            accessToken: PUBLIC_MAPBOX_ACCESS_TOKEN,
            container: "map-container",
            style: "mapbox://styles/mapbox/streets-v11",
            center: [19.21, 52.11],
            zoom: 5,
        });
        currentMap.set(map);
    });
</script>

<div id="map-container"></div>
<style>
    #map-container { position: absolute; top: 0; bottom: 0; width: 100%; }
</style>