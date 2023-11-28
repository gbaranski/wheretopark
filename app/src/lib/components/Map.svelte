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
	import ParkingMapMarker from './ParkingMapMarker.svelte';
	import 'mapbox-gl/dist/mapbox-gl.css';
	import { SpotType, type ID, ParkingLot, Feature } from '$lib/parkingLot';
	import { geolocation } from '$lib/utils';

	onMount(async () => {
		const mapboxgl = await import('mapbox-gl');
		const map = new mapboxgl.Map({
			accessToken: PUBLIC_MAPBOX_ACCESS_TOKEN,
			container: 'map-container',
			style: 'mapbox://styles/mapbox/navigation-day-v1',
			center: [19.21, 52.11],
			zoom: 5
		});
		const geolocateControl = new mapboxgl.GeolocateControl({
			positionOptions: {
				enableHighAccuracy: true
			},
			// When active the map will receive updates to the device's location as it changes.
			trackUserLocation: true,
			// Draw an arrow next to the location dot to indicate which direction the device is heading.
			showUserHeading: true,
			geolocation: geolocation()
		});
		map.addControl(geolocateControl, 'bottom-right');
		currentMap.set(map);
		const currentMarkers: mapboxgl.Marker[] = [];
		const updateMap = (parkingLots: Record<ID, ParkingLot>, filters: SearchFilters) => {
			// console.log('updating map', { parkingLots, filters });
			currentMarkers.forEach((marker) => {
				marker.remove();
			});
			Object.entries(parkingLots)
				.filter(([_, parkingLot]) => {
					const status = parkingLot.status(SpotType.car);
					if (filters.openNow && !status.isOpen()) return false;

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
					new ParkingMapMarker({ target: markerElement, props: { parkingLot } });
					const marker = new mapboxgl.Marker(markerElement).setLngLat([longitude, latitude]);
					marker.getElement().addEventListener('click', () => {
						goto(`/parking-lot/${id}`, { noScroll: true });
					});
					return marker;
				})
				.forEach((marker) => {
					marker.addTo(map);
					currentMarkers.push(marker);
				});
			map.resize();
		};
		parkingLotsStore.subscribe((parkingLots) => updateMap(parkingLots, $searchFilters));
		searchFilters.subscribe((filters) => updateMap($parkingLotsStore, filters));
	});
</script>

<div id="map-container" class="w-full h-full absolute" />
