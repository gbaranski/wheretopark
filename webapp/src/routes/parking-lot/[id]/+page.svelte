<script lang="ts">
	import { currentMap } from "$lib/store";
	import { SpotType, type ParkingLot, Feature, type State, type Metadata } from "$lib/types";
	import { capitalizeFirstLetter, getCategory, googleMapsLink, humanizeDuration, resourceIcon, resourceText, timeFromNow } from "$lib/utils";
    import { Title, Text, Divider } from '@svelteuidev/core';
    import { page } from '$app/stores';

    export let data: {parkingLot: ParkingLot};

    $: metadata = data.parkingLot.metadata as Metadata;
    $: state = data.parkingLot.state as State;
    $: features = metadata?.features?.map((feature) => Feature[feature as keyof typeof Feature]);
    $: category = getCategory(features || []);
    
    $: {
        const [longitude, latitude] = metadata.geometry.coordinates;
        $currentMap?.flyTo({
            center: [longitude, latitude],
            zoom: 15
        });
    }
</script>

<svelte:head>
    <title>{metadata.name}</title>
	<meta name="description" content="Details of {capitalizeFirstLetter(category)} parking lot in {metadata.name} at {metadata.address}, containing prices, opening hours and it's availability of parking spots."/>
	<meta name="keywords" content="{metadata.name}, {metadata.address}, Parking Lot, Smart City, Gdańsk, Gdynia, Sopot, Tricity"/>
</svelte:head>
<div class="container">
    <span style="display: flex;">
        <Title root={"span"} size={26} override={{flex: 1}}>{metadata.name}</Title>
        <a style="text-align: right;" href={googleMapsLink(metadata.geometry)} target="_blank" rel="noreferrer">
            <i class="material-icons">directions</i>
        </a>
    </span>
    <Text size={14} weight={"semibold"}>{category}</Text>
    <Divider/>
    
    <div class="field">
        <i class="material-icons">place</i>
        <Text size={14} root="span" >
            {metadata.address}
        </Text>
    </div>

    <div class="field">
        <i class="material-icons">directions_car</i>
        <Text size={14} root="span" >
            {state.availableSpots[SpotType[SpotType.CAR]]} available spots
        </Text>
    </div>
    
    <div class="field">
        <i class="material-icons">schedule</i>
        <Text size={14} root="span" >
            Last updated {timeFromNow(state.lastUpdated)}
        </Text>
    </div>
    
    {#each metadata.resources as resource}
        <div class="field">
                <i class="material-icons">{resourceIcon(resource)}</i>
                <a href={resource}>
                    <Text size={14} root="span" >
                        {resourceText(resource)}
                    </Text>
                </a>
        </div>
    {/each}
    {#each metadata.rules as rule}
        <div style="margin-top: 20px;">
            <Text weight="semibold">{rule.hours}</Text>
            {#each rule.pricing as pricing}
                <Text weight="light" size={16}>{pricing.repeating ? "Each " : ""}{humanizeDuration(pricing.duration)} - {pricing.price}{metadata.currency}</Text>
            {/each}
        </div>
    {/each}
</div>
<style>
    .container {
        padding: 20px;
    }
    
    .field {
        margin-bottom: 10px;
        
    }
    
    .field > i {
        font-size: 18px !important;
        margin-right: 10px;
        vertical-align: middle;
    }
    
    div.field > :global(*) {
        font-weight: 500;
        vertical-align: middle;
    }
    
</style>