<script lang="ts">
	import { MAPBOX_ACCESS_TOKEN } from "$lib/environment";
	import { onMount } from "svelte";
    import 'mapbox-gl/dist/mapbox-gl.css';
	import { SpotType, type ID, type ParkingLot } from "$lib/types";
	import { currentMap } from "$lib/store";
	import { markerColor, parkingLotStatus } from "$lib/utils";
	import MapMarker from "./MapMarker.svelte";
    
    export let parkingLots: Record<ID, ParkingLot>;

    const markers: Record<ID, MapMarker> = {};

    onMount(async () => {
        const mapboxgl = await import("mapbox-gl");
        const map = new mapboxgl.Map({
            accessToken: MAPBOX_ACCESS_TOKEN,
            container: "map-container",
            style: "mapbox://styles/mapbox/streets-v11",
            center: [18.64, 54.35],
            zoom: 10,
        });
        currentMap.set(map);
        
        Object.entries(parkingLots).map(([id, parkingLot]) => {
            const [longitude, latitude] = parkingLot.metadata.geometry.coordinates;
            const popupHtml = `
                <h4>${parkingLot.metadata.name}</h4>
                <p>${parkingLot.state.availableSpots["CAR"]} available car spots</p>
                <a href="/parking-lot/${id}">Open</a>
            `;
            const popup = new mapboxgl.Popup({offset: 25}).setHTML(popupHtml);
            const status = parkingLotStatus(parkingLot, SpotType.CAR)[0];
            const options: mapboxgl.MarkerOptions = {
                color: markerColor(parkingLot.state.availableSpots["CAR"], parkingLot.metadata.totalSpots["CAR"], status)
            }
            const marker = markers[id];
            return new mapboxgl.Marker(options)
                .setLngLat([longitude, latitude])
                .setPopup(popup);
        }).forEach((marker) => marker.addTo(map));

    });
</script>

<div id="map-container"></div>
{#each Object.entries(parkingLots) as [id, parkingLot]}
    <MapMarker parkingLot={parkingLot} bind:this={markers[id]} />
{/each}
<style>
    #map-container { position: absolute; top: 0; bottom: 0; width: 100%; }
</style>