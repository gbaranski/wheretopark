<script lang="ts">
	import { PUBLIC_MAPBOX_ACCESS_TOKEN } from '$env/static/public';
	import { onMount } from 'svelte';
	import {
		currentMap,
		parkingLots as parkingLotsStore,
		searchFilters,
		type SearchFilters
	} from '$lib/store';
	import { goto } from '$app/navigation';
	import MapMarker from './MapMarker.svelte';
	import 'mapbox-gl/dist/mapbox-gl.css';
	import { SpotType, type ID, ParkingLot, Status, Feature } from '$lib/parkingLot';

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
		const currentMarkers: mapboxgl.Marker[] = [];
		const updateMap = (parkingLots: Record<ID, ParkingLot>, filters: SearchFilters) => {
			console.log('updating map', { parkingLots, filters });
			currentMarkers.forEach((marker) => {
				marker.remove();
			});
			Object.entries(parkingLots)
				.filter(([_, parkingLot]) => {
					const status = parkingLot.status(SpotType.car);
					if (filters.openNow && (status.isOpen() || status.isOpeningSoon()))
						return false;

					const requiredFeatures = Object.entries(filters.hasFeatures).filter(
						([_, required]) => required
					);
					if (requiredFeatures.length > 0) {
						const hasFeature = requiredFeatures.some(([featureStr, _]) => {
							const feature = featureStr as keyof typeof Feature;
							return parkingLot.features.includes(Feature[feature]);
						});
						if (!hasFeature) return false;
					}
					if (filters.minAvailableSpots > parkingLot.availableSpotsFor(SpotType.car)) return false;
					return true;
				})
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
				.forEach((marker) => {
					marker.addTo(map);
					currentMarkers.push(marker);
				});
		};
		parkingLotsStore.subscribe((parkingLots) => updateMap(parkingLots, $searchFilters));
		searchFilters.subscribe((filters) => updateMap($parkingLotsStore, filters));
	});
</script>

<div
	id="map-container"
	class="w-full h-4/6 absolute top-16 lg:right-0 lg:h-[calc(100%-4rem)] lg:w-7/12 xl:w-8/12 2xl:w-9/12"
></div>
