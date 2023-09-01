<script lang="ts">
	import { PUBLIC_MAPBOX_ACCESS_TOKEN } from '$env/static/public';
	import { onMount } from 'svelte';
	import { currentMap, parkingLots as parkingLotsStore } from '$lib/store';
	import { goto } from '$app/navigation';
	import MapMarker from './MapMarker.svelte';
	import 'mapbox-gl/dist/mapbox-gl.css';


	onMount(async () => {
		const mapboxgl = await import('mapbox-gl');
		const map = new mapboxgl.Map({
			accessToken: PUBLIC_MAPBOX_ACCESS_TOKEN,
			container: 'map-container',
			style: 'mapbox://styles/mapbox/navigation-day-v1',
			center: [19.21, 52.11],
			zoom: 5
		});
		currentMap.set(map);
		parkingLotsStore.subscribe((parkingLots) => {
			Object.entries(parkingLots)
				.map(([id, parkingLot]) => {
					const [longitude, latitude] = parkingLot.geometry.coordinates;
					const markerElement = document.createElement('div');
					new MapMarker({  target: markerElement, props: {parkingLot} });
					const marker = new mapboxgl.Marker(markerElement).setLngLat([longitude, latitude]);
					marker.getElement().addEventListener('click', () => {
						goto(`/parking-lot/${id}`);
					})
					return marker;
				})
				.forEach((marker) => marker.addTo(map));
		});
	});
</script>

<div id="map-container"></div>

<style>
	#map-container {
		position: absolute;
		top: 0;
		bottom: 0;
		width: 100%;
	}
</style>
