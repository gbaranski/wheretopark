<script lang="ts">
	import { MAPBOX_ACCESS_TOKEN } from "$lib/environment";
	import mapboxgl from "mapbox-gl";
	import { onMount } from "svelte";
    import 'mapbox-gl/dist/mapbox-gl.css';
	import type { ID, ParkingLot } from "$lib/types";
    
    export let parkingLots: Record<ID, ParkingLot>;

    onMount(() => {
        mapboxgl.accessToken = MAPBOX_ACCESS_TOKEN;
        const map = new mapboxgl.Map({
            container: "map-container",
            style: "mapbox://styles/mapbox/streets-v11",
            center: [18.64, 54.35],
            zoom: 10,
        });
        
        Object.entries(parkingLots).map(([id, parkingLot]) => {
            const [longitude, latitude] = parkingLot.metadata.location.geometry.coordinates;
            return new mapboxgl.Marker()
                .setLngLat([longitude, latitude]);
        }).forEach((marker) => marker.addTo(map));

    });
</script>

<div id="map-container"></div>
<style>
    #map-container { position: absolute; top: 0; bottom: 0; width: 100%; }
</style>