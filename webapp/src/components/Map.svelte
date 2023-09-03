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
		// Add geolocate control to the map.
		map.addControl(
			new mapboxgl.GeolocateControl({
				positionOptions: {
					enableHighAccuracy: true
				},
				// When active the map will receive updates to the device's location as it changes.
				trackUserLocation: true,
				// Draw an arrow next to the location dot to indicate which direction the device is heading.
				showUserHeading: true
			})
		);
		currentMap.set(map);
		parkingLotsStore.subscribe((parkingLots) => {
			Object.entries(parkingLots)
				.map(([id, parkingLot]) => {
					const [longitude, latitude] = parkingLot.geometry.coordinates;
					const markerElement = document.createElement('div');
					new MapMarker({ target: markerElement, props: { parkingLot } });
					const marker = new mapboxgl.Marker(markerElement).setLngLat([longitude, latitude]);
					marker.getElement().addEventListener('click', () => {
						goto(`/parking-lot/${id}`);
					});
					return marker;
				})
				.forEach((marker) => marker.addTo(map));
		});
	});
</script>

<div
	id="map-container"
	class="w-full h-4/6 absolute top-16 lg:right-0 lg:h-[calc(100%-4rem)] lg:w-7/12 xl:w-8/12 2xl:w-9/12"
></div>
