<script lang="ts">
	import { LL } from '$lib/i18n/i18n-svelte';
	import { awards } from '$lib/assets/award';
	import { cities } from '$lib/assets/city';
	import { news } from '$lib/assets/news';
	import appPreview from '$lib/assets/preview.webp';
	import exampleParkingLot from '$lib/assets/parking-basen-klodzko.webp';
	import { getRandomBetween, mailFromOperator } from '$lib/utils';
	import { inview } from 'svelte-inview';
	import AppStoreBadge from '$lib/assets/badge/AppStore.svelte';
	import GooglePlayBadge from '$lib/assets/badge/GooglePlay.svelte';
	import { iosScreenshot } from '$lib/assets/screenshots';

	let counter = 0;
	const onCounterVisible = () => {
		if (counter > 0) return;
		const interval = setInterval(() => {
			counter = getRandomBetween(0, 100);
		}, 10);
		setTimeout(() => {
			clearInterval(interval);
			counter = 10;
		}, 2000);
	};

	const references = {
		krakow2018:
			'https://yadda.icm.edu.pl/baztech/element/bwmeta1.element.baztech-51e98198-f84e-4b3f-8a4a-2492f1fef1df',
		krakow2021:
			'https://yadda.icm.edu.pl/baztech/element/bwmeta1.element.baztech-d8bf4176-2ea0-4571-b5a5-bdc066343c33',
		copark:
			'https://www.researchgate.net/publication/307534127_Investigating_Vehicle-to-Vehicle_Communication_for_Cooperative_Car_Parking_the_CoPark_Approach'
	};
</script>


<div class="flex flex-col w-full gap-3">
	<div class="flex flex-col gap-12 md:gap-0 md:flex-row justify-between">
		<div class="flex flex-col justify-center gap-8">
			<div class="flex flex-col gap-3">
				<h1><b>{$LL.heading()}</b></h1>
				<h2 class="font-regular text-lg lg:text-xl pt-2">
					{@html $LL.subheading()} 🚗.
				</h2>
			</div>
			<div class="max-md:hidden py-2" />
			<div class="flex flex-col justify-center">
				<div class="flex flex-row flex-wrap items-center justify-center">
					<AppStoreBadge />
					<GooglePlayBadge />
				</div>
				<div class="text-center">
					<span>{$LL.or()}</span>
					<a class="link link-info" href="https://web.wheretopark.app">{$LL.openAppInBrowser()}</a>
				</div>
			</div>
		</div>

		<div class="self-center justify-center">
			<img src={iosScreenshot} alt="ios app screenshot" width="256px" />
		</div>
	</div>

	<div class="self-center divider w-96 p-12" />
	<div class="flex justify-evenly items-center flex-col">
		<h2 class="font-extrabold text-3xl pb-5">{$LL.problem.what()}</h2>
		<div class="flex justify-evenly flex-row flex-wrap gap-10">
			<div>
				<h2 class="font-extrabold text-5xl text-info">
					<span use:inview on:inview_enter={onCounterVisible}>
						{$LL.minutes(counter)}
					</span>
				</h2>
				<p>
					{@html $LL.AVERAGE_TIME_WASTED()}
				</p>
			</div>
			<div>
				<h2 class="font-extrabold text-5xl text-warning">93.7%</h2>
				<p>{@html $LL.OF_DRIVERS_CIRCULATE()}</p>
			</div>
			<div>
				<h2 class="font-extrabold text-5xl text-error">
					{@html $LL.tonsOf('CO<sub>2</sub>', 730)}
				</h2>
				<p>
					{@html $LL.AIR_POLUTION_ANNUAL_CONTRIBUTION('Los Angeles')}
				</p>
			</div>
		</div>
	</div>

	<div class="self-center divider w-96 p-12" />
	<div class="flex justify-evenly items-center flex-col lg:flex-row-reverse gap-8">
		<div class="">
			<h2 class="font-extrabold text-4xl">{$LL.OUR_SOLUTION()}</h2>
			<h3 class="text-md">{$LL.OUR_SOLUTION_TEXT()}</h3>
		</div>
		<img
			class="w-96 lg:w-2/3 rounded-2xl"
			src={appPreview}
			alt="preview of the app"
			width="384px"
			height="207px"
		/>
	</div>

	<div class="self-center divider w-96 p-12" />
	<div class="flex justify-evenly items-center flex-col-reverse lg:flex-row gap-8">
		<div class="">
			<h2 class="font-extrabold text-3xl">{$LL.POWERED_BY_AI()} 🤖</h2>
			<h3 class="text-md">
				{@html $LL.OUR_SYSTEM_MARKS()}
			</h3>
		</div>
		<img
			class="w-96 lg:w-2/3 rounded-lg"
			src={exampleParkingLot}
			alt="animation of our ai"
			width="384px"
			height="216px"
		/>
	</div>

	<div class="self-center divider w-96 p-12" />
	<h2 class="font-extrabold text-4xl">{$LL.WE_HAVE_PARKNG_LOTS()}</h2>
	<div class="pt-10 flex flex-row flex-wrap gap-16 justify-center">
		{#each cities as city}
			<img
				src={city.image}
				alt="{city} logo"
				class="w-32 object-scale-down"
				width="128px"
				height="128px"
			/>
		{/each}
	</div>

	<div class="self-center divider w-96 p-12" />
	<h2 class="font-extrabold text-3xl">{$LL.AWARDS()} &nbsp;🏆</h2>
	<div class="pt-10 flex flex-row flex-wrap gap-16 justify-center">
		{#each awards as award}
			<a class="w-32" href={award.link} target="_blank">
				<img
					src={award.image}
					alt="{new URL(award.link).hostname} logo"
					class="h-32 object-scale-down pb-2"
					width="128px"
					height="128px"
				/>
				<span class="link link-secondary">{award.description}</span>
			</a>
		{/each}
	</div>

	<div class="self-center divider w-96 p-12" />
	<h2 class="font-extrabold text-3xl">{$LL.TALKS_ABOUT_US()} &nbsp;📢</h2>
	<div class="pt-10 flex flex-row flex-wrap gap-16 justify-center">
		{#each news as entry}
			<a class="w-32" href={entry.link} target="_blank">
				<img
					src={entry.image}
					alt="{new URL(entry.link).hostname} logo"
					class="h-32 object-scale-down"
					height="128px"
					width="128px"
				/>
			</a>
		{/each}
	</div>

	<div class="self-center divider w-96 p-12" />
	<div>
		<div class="pb-10">
			<h2 class="font-extrabold text-3xl pb-10">{$LL.BENEFITS()} &nbsp;📈</h2>
			<div
				class="flex flex-col text-left items-center justify-evenly gap-10 lg:flex-row lg:items-start"
			>
				<div class="flex-1 w-3/4">
					<h2 class="text-2xl font-bold">{$LL.FOR_DRIVERS()} 🚗</h2>
					<ul class="text-sm list-disc">
						<li>
							{@html $LL.FOR_DRIVERS_SAVINGS(
								`<a class="link" href="${references.krakow2021}" target="_blank">2$</a>`
							)}
						</li>
						<li>
							{@html $LL.FOR_DRIVERS_DISABILITIES()}
						</li>
						<li>
							{@html $LL.FOR_DRIVERS_PREDICTIONS()}
						</li>
					</ul>
				</div>
				<div class="flex-1 w-3/4">
					<h2 class="text-2xl font-bold">{$LL.FOR_OPERATORS()} 👔</h2>
					<ul class="text-sm list-disc">
						<li>
							{@html $LL.FOR_OPERATORS_PROMOTION()}
						</li>
						<li>
							{@html $LL.FOR_OPERATORS_RESERVATIONS()}
						</li>
					</ul>
				</div>
				<div class="flex-1 w-3/4">
					<h2 class="text-2xl font-bold">{$LL.FOR_CITIES()} 🌇</h2>
					<ul class="text-sm list-disc">
						<li>
							{@html $LL.FOR_CITIES_ANALYSIS()}
						</li>
						<li>
							{@html $LL.FOR_CITIES_INCORRECTLY_PARKED()}
						</li>
						<li>
							{@html $LL.FOR_CITIES_DECREASED_TRAFFIC(
								`<a href="${references.krakow2021}" class="link" target="_blank">70%</a>`
							)}
						</li>
						<li>
							{@html $LL.FOR_CITIES_OPTIMISED_DISTRIBUTION()}
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>

	<div class="self-center divider w-96 p-12" />
	<div>
		<div class="pb-10">
			<h2 class="font-extrabold text-4xl">{$LL.ARE_YOU_PARKING_OPERATOR()}</h2>
		</div>
		<a href={mailFromOperator} class="btn btn-primary">{$LL.CONTACT_US()}</a>
	</div>

	<div class="pb-12" />
</div>
