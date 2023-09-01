<script lang="ts">
	import { currentMap } from '$lib/store';
	import { SpotType, type ParkingLot, Feature, type State, type Metadata } from '$lib/types';
	import {
		availabilityColor,
		capitalizeFirstLetter,
		getCategory,
		getWeekday,
		googleMapsLink,
		humanizeDuration,
		parkingLotStatus,
		statusColor,
		preferredComment,
		resourceText,
		rulesForDay,
		timeFromNow,
		weekdays
	} from '$lib/utils';
	import Markdown from 'svelte-markdown';
	import ResourceIcon from '$components/ResourceIcon.svelte';

	export let data: { parkingLot: ParkingLot };
	$: metadata = data.parkingLot.metadata;
	$: state = data.parkingLot.state;
	$: [status, statusComment] = parkingLotStatus(data.parkingLot, SpotType.CAR);
	$: features = metadata.features.map((feature) => Feature[feature as keyof typeof Feature]);
	$: category = getCategory(features || []);
	$: comment = preferredComment(metadata.comment || {});
	$: avColor = availabilityColor(state.availableSpots['CAR'], metadata.totalSpots['CAR']);
	$: stColor = statusColor(status);

	$: {
		const [longitude, latitude] = metadata.geometry.coordinates;
		$currentMap?.flyTo({
			center: [longitude, latitude],
			zoom: 15
		});
	}

	let selectedWeekday = getWeekday();
	$: applicableRules = rulesForDay(metadata.rules, SpotType.CAR, selectedWeekday);

	const onShare = async () => {
		if (isSharing) {
			isSharing = false;
			return;
		};
		try {
			if (navigator.canShare != null && navigator.canShare()) {
				await navigator.share({
					url: window.location.href,
					title: metadata.name,
					text: `Check out ${metadata.name} parking lot in wheretopark.app!`
				});
			} else {
				console.log("copied to clipboard")
				await navigator.clipboard.writeText(window.location.href);
				isSharing = true;
			}
		} catch (error) {
			console.log(error);
		}
	};

	let isSharing = false;
</script>

<svelte:head>
	<title>Parking {metadata.name.replace('Parking', '')}</title>
	<meta
		name="description"
		content="Details of {capitalizeFirstLetter(
			category
		)} parking lot in {metadata.name} at {metadata.address}, containing prices, opening hours and it's availability of parking spots."
	/>
	<meta
		name="keywords"
		content="{metadata.name}, {metadata.address}, Parking Lot, Smart City, Gdańsk, Gdynia, Sopot, Tricity"
	/>
</svelte:head>
<div class="pt-10 pl-5 w-11/12">
	<h1 class="font-sans text-3xl font-extrabold">{metadata.name}</h1>
	<h2 class="font-mono text-sm font-light mb-2">{category}</h2>
	<div class="join w-full">
		<a
			class="btn btn-primary rounded-md w-2/3"
			href={googleMapsLink(metadata.geometry)}
			target="_blank"
		>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				class="w-6 h-6"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M8.25 18.75a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m3 0h6m-9 0H3.375a1.125 1.125 0 01-1.125-1.125V14.25m17.25 4.5a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m3 0h1.125c.621 0 1.129-.504 1.09-1.124a17.902 17.902 0 00-3.213-9.193 2.056 2.056 0 00-1.58-.86H14.25M16.5 18.75h-2.25m0-11.177v-.958c0-.568-.422-1.048-.987-1.106a48.554 48.554 0 00-10.026 0 1.106 1.106 0 00-.987 1.106v7.635m12-6.677v6.677m0 4.5v-4.5m0 0h-12"
				/>
			</svg>
			Navigate
		</a>

		<div class="dropdown dropdown-end w-full ml-5">
			<!-- svelte-ignore a11y-no-noninteractive-tabindex -->
			<!-- svelte-ignore a11y-label-has-associated-control -->
			<label tabindex="0" class="btn btn-neutral rounded-md w-full">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="w-6 h-6"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M6.75 12a.75.75 0 11-1.5 0 .75.75 0 011.5 0zM12.75 12a.75.75 0 11-1.5 0 .75.75 0 011.5 0zM18.75 12a.75.75 0 11-1.5 0 .75.75 0 011.5 0z"
					/>
				</svg>
			</label>
			<!-- svelte-ignore a11y-no-noninteractive-tabindex -->
			<ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
				<li>
					<label class="swap justify-start">
						<input type="checkbox" value={isSharing} on:change={onShare} />
						<span class="swap-off inline-flex gap-2">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="1.5"
								stroke="currentColor"
								class="w-5 h-5"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M7.217 10.907a2.25 2.25 0 100 2.186m0-2.186c.18.324.283.696.283 1.093s-.103.77-.283 1.093m0-2.186l9.566-5.314m-9.566 7.5l9.566 5.314m0 0a2.25 2.25 0 103.935 2.186 2.25 2.25 0 00-3.935-2.186zm0-12.814a2.25 2.25 0 103.933-2.185 2.25 2.25 0 00-3.933 2.185z"
								/>
							</svg>
							<span>Share</span>
						</span>
						<div class="swap-on">Copied to clipboard</div>
					</label>
				</li>
			</ul>
		</div>
	</div>
	<div class="stats w-full ml-0 left-0 p-0">
		<div class="stat">
			<div class="stat-title text-xs font-mono font-semibold">AVAILABLE SPACES</div>
			<div class="stat-value font-mono font-extrabold" style="color:{avColor}">
				{state.availableSpots[SpotType[SpotType.CAR]]}
			</div>
			<div class="stat-desc">
				Updated <span class="text-success">{timeFromNow(state.lastUpdated)}</span>
			</div>
		</div>
		<div class="stat">
			<div class="stat-title text-xs font-mono font-semibold">CURRENT STATUS</div>
			<div class="stat-value font-mono" style="color:{stColor}">{status}</div>
			<div class="stat-desc">{statusComment}</div>
		</div>
	</div>

	<div class="divider"></div>

	<div class="flex flex-auto flex-row justify-between">
		<div>
			{#each applicableRules as rule}
				<p class="text-xl font-bold mt-5">{rule.humanHours}</p>
				{#each rule.pricing as pricing}
					<div>
						{pricing.repeating ? 'Each ' : ''}
						{humanizeDuration(pricing.duration)} - {pricing.price}
						{metadata.currency}
					</div>
				{/each}
			{/each}
		</div>

		<select class="select select-bordered max-w-xs mt-4" bind:value={selectedWeekday}>
			<option disabled>Pick a day</option>
			{#each weekdays as weekday, i}
				<option value={i}>{weekday}</option>
			{/each}
		</select>
	</div>

	<div class="divider"></div>

	<h2 class="text-2xl font-bold mb-3">Additional info</h2>
	<div class="flex flex-col gap-y-2 font-mono text-sm">
		<p class="inline-flex">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				class="w-5 h-5"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z"
				/>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z"
				/>
			</svg>
			<span class="ml-2">
				{metadata.address}
			</span>
		</p>

		{#if (metadata.paymentMethods?.length || 0) > 0}
			<p class="inline-flex">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="w-6 h-6"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M2.25 8.25h19.5M2.25 9h19.5m-16.5 5.25h6m-6 2.25h3m-3.75 3h15a2.25 2.25 0 002.25-2.25V6.75A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25v10.5A2.25 2.25 0 004.5 19.5z"
					/>
				</svg>

				{#each metadata.paymentMethods as paymentMethod}
					<div class="badge badge-neutral badge-lg ml-2 font-mono text-xs">{paymentMethod}</div>
				{/each}
			</p>
		{/if}

		{#each metadata.resources as resource}
			{@const url = new URL(resource)}
			<div class="inline-flex">
				<ResourceIcon resource={url} />
				<a class="link ml-2" href={url.toString()}>
					{resourceText(url)}
				</a>
			</div>
		{/each}
	</div>

	{#if comment}
		<div class="divider"></div>
		<h2 class="text-2xl font-bold mb-3">Description</h2>
		<article class="prose prose-sm">
			<Markdown source={comment} />
		</article>
	{/if}
</div>
