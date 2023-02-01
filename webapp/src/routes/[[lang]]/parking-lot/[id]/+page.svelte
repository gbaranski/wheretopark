<script lang="ts">
	import { currentMap } from "$lib/store";
	import { SpotType, type ParkingLot, Feature, type State, type Metadata } from "$lib/types";
	import { capitalizeFirstLetter, getCategory, googleMapsLink, humanizeDuration, parkingLotStatus, parkingLotStatusColor, preferredComment, resourceIcon, resourceText, spotTypeIcon, timeFromNow } from "$lib/utils";
    import { Title, Text, Divider, Tooltip, Anchor, Badge } from '@svelteuidev/core';
    import { LL } from '$lib/i18n/i18n-svelte';
    import Markdown from "svelte-markdown";

    export let data: {parkingLot: ParkingLot};
    const [status, statusComment] = parkingLotStatus(data.parkingLot, SpotType.CAR);

    $: metadata = data.parkingLot.metadata as Metadata;
    $: state = data.parkingLot.state as State;
    $: features = metadata?.features?.map((feature) => Feature[feature as keyof typeof Feature]);
    $: category = getCategory(features || []);
    $: comment = preferredComment(metadata.comment || {});
    
    $: {
        const [longitude, latitude] = metadata.geometry.coordinates;
        $currentMap?.flyTo({
            center: [longitude, latitude],
            zoom: 15
        });
    }
</script>

<svelte:head>
    <title>Parking {metadata.name.replace("Parking", "")}</title>
	<meta name="description" content="Details of {capitalizeFirstLetter(category)} parking lot in {metadata.name} at {metadata.address}, containing prices, opening hours and it's availability of parking spots."/>
	<meta name="keywords" content="{metadata.name}, {metadata.address}, Parking Lot, Smart City, Gdańsk, Gdynia, Sopot, Tricity"/>
</svelte:head>
<div class="container">
    <span style="display: flex;">
        <Text root={"span"} size={26} override={{flex: 1}}>{metadata.name}</Text>
        <Anchor root="a" external override={{textAlign: 'right'}} href={googleMapsLink(metadata.geometry)}>
            <i class="material-icons">directions</i>
        </Anchor>
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
        <Tooltip label="Last updated {timeFromNow(state.lastUpdated)}">
            <Text size={14} root="span" weight={400}>
                {$LL.AVAILABLE_SPOTS({count: state.availableSpots[SpotType[SpotType.CAR]]})}
            </Text>
        </Tooltip>
    </div>
    
    <div class="field">
        <i class="material-icons">schedule</i>
        <Text size={14} root="span" override={{color: parkingLotStatusColor(status)}}>
            {status}
        </Text>
        {#if statusComment != undefined}
            <Text size={14} root="span">
            -  {statusComment}
            </Text>
        {/if}
    </div>

    {#if (metadata.paymentMethods?.length || 0) > 0}
        <div class="field">
            <i class="material-icons">payment</i>
                
            {#each metadata.paymentMethods as paymentMethod}
                <span style="margin-right: 1px;">
                    <Badge size="sm" variant="outline" >
                        {paymentMethod}
                    </Badge>
                </span>
            {/each}
        </div>
    {/if}

    {#each metadata.resources as resource}
        <div class="field">
                <i class="material-icons">{resourceIcon(resource)}</i>
                <Anchor size={14} root="a" external href={resource} color="inherit">
                        {resourceText(resource)}
                </Anchor>
        </div>
    {/each}

    <div class="field">
        <i class="material-icons" style="position: absolute;">paid</i>
        {#each metadata.rules as rule, i}
            <div style="margin-left: 32px; margin-top: 10px;">
                <Text root="span" weight="semibold">{rule.hours}</Text>
                {#each rule.applies || [] as spotType}
                    <i class="material-icons" style="float: right; font-size: 18px;">{spotTypeIcon(spotType)}</i>
                {/each}

                <div style="margin-left: 10px;">
                    {#each rule.pricing as pricing}
                        <Text weight="light" size={16}>{pricing.repeating ? "Each " : ""}{humanizeDuration(pricing.duration)} - {pricing.price}{metadata.currency}</Text>
                    {/each}
                </div>
            </div>
        {/each}
    </div>
    

    <Divider variant="dashed" />
    {#if comment}
        <Text weight="light" size={14}>
            <Markdown source={comment} />
        </Text>
    {/if}
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