<script lang="ts">
	import { SpotType, type ParkingLot } from '$lib/parkingLot';
	import MapMarker from './MapMarker.svelte';

	export let parkingLot: ParkingLot;
	const status = parkingLot.status(SpotType.car);
	const rating = parkingLot.rating(status, SpotType.car)
</script>

<MapMarker rating={rating} />