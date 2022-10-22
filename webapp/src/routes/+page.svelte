<script lang="ts">
	import List from "../components/List.svelte";
	import type { LayoutData } from "../types/layout";
	
    export let data: LayoutData;
</script>

<List parkingLots={data.parkingLots}/>