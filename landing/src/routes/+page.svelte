<script lang="ts">
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
		}, 3000);
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

	<div class="flex flex-row justify-between flex-wrap items-center">
		<div class="flex flex-col justify-center gap-8">
			<div class="flex flex-col gap-3">
				<h1 class="font-bold text-4xl">Your parking lot assistant.</h1>
				<h2 class="font-regular text-lg lg:text-xl pt-2">
					with the help of <b>AI</b> we'll find you an <b>available spot</b> nearby 🚗.
				</h2>
			</div>
			<div class="py-2"/>
			<div class="flex flex-col justify-center">
				<div class="flex flex-row flex-wrap items-center justify-center">
					<AppStoreBadge />
					<GooglePlayBadge />
				</div>
				<div class="text-center">
					<span>or</span>
					<a class="link link-info" href="https://web.wheretopark.app">open app in the browser</a>
				</div>
			</div>
		</div>

		<div>
			<img src={iosScreenshot} alt="ios app screenshot" width="256px" />
		</div>
	</div>

	<div class="self-center divider w-96 p-12" />
	<div class="flex justify-evenly items-center flex-col">
		<h2 class="font-extrabold text-3xl pb-5">What problem are we solving?</h2>
		<div class="flex justify-evenly flex-row flex-wrap gap-10">
			<div>
				<h2 class="font-extrabold text-5xl text-info">
					<span use:inview on:inview_enter={onCounterVisible}>{counter}</span> minutes
				</h2>
				<p>
					Average <span class="font-bold">time wasted</span> each time you want to park.
				</p>
			</div>
			<div>
				<h2 class="font-extrabold text-5xl text-warning">93.7%</h2>
				<p>
					Of drivers <b>circulate around</b> parking lots
				</p>
			</div>
			<div>
				<h2 class="font-extrabold text-5xl text-error">
					730 tons of CO<sub>2</sub>
				</h2>
				<p>
					Annual contribution to <b>air pollution</b> in Los Angeles alone.
				</p>
			</div>
		</div>
	</div>

	<div class="self-center divider w-96 p-12" />
	<div class="flex justify-evenly items-center flex-col lg:flex-row-reverse gap-8">
		<div class="">
			<h2 class="font-extrabold text-5xl">Our solution</h2>
			<h3 class="text-xl">Easy access to information about available parking lots nearby</h3>
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
			<h2 class="font-extrabold text-4xl">Powered by AI &nbsp;🤖</h2>
			<h3 class="text-xl">
				Our system automatically marks
				<span class="text-green-600">green</span>
				free spaces, and
				<span class="text-red-600">red</span>
				occupied.
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
	<h2 class="font-extrabold text-4xl">We have parkings lots in</h2>
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
	<h2 class="font-extrabold text-4xl">Awards &nbsp;🏆</h2>
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
	<h2 class="font-extrabold text-4xl">Talks about us &nbsp;📢</h2>
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
			<h2 class="font-extrabold text-4xl pb-10">Benefits 📈</h2>
			<div
				class="flex flex-col text-left items-center justify-evenly gap-10 lg:flex-row lg:items-start"
			>
				<div class="flex-1 w-3/4">
					<h2 class="text-2xl font-bold">For drivers 🚗</h2>
					<ul class="list-disc">
						<li>
							<b>Money savings</b> of around
							<a class="link" href={references.krakow2021} target="_blank">9zł or 2$</a> per each parking
							session (based on costs of running the car and wasted time).
						</li>
						<li>
							<b>People with disabilities</b> can easily check which parking lots have specially designated
							spaces for them, and whether these spaces are free.
						</li>
						<li>
							<b>Predictions of availability</b>. The user has access to the statistics of the
							occupancy of places on parking lots at a given time and on a given day. For example,
							on Friday at 5 pm, the parking lot is usually 80% full, and on Saturday at 2 pm, it is
							50% full.
						</li>
					</ul>
				</div>
				<div class="flex-1 w-3/4">
					<h2 class="text-2xl font-bold">For parking operators 👔</h2>
					<ul class="list-disc">
						<li>
							<b>In-app promotion</b> - more customers.
						</li>
						<li>
							<b>Reserving parking spaces</b> directly in the app, without need of phone calling.
						</li>
					</ul>
				</div>
				<div class="flex-1 w-3/4">
					<h2 class="text-2xl font-bold">For cities 🌇</h2>
					<ul class="list-disc">
						<li>
							<b>Trend analysis</b> - will help cities choosing the right place for a new parking lot,
							or managing the prices.
						</li>
						<li>
							Our system can <b>detect illegally parked vehicle</b> and notify the authorities.
						</li>
						<li>
							<b>Decreased car traffic</b> - in some cities up to
							<a href={references.krakow2021} class="link" target="_blank">70% of traffic</a> is caused
							by drivers.
						</li>
						<li>
							<b>Optimised distribution of cars</b> - drivers could start using less known, or recently
							opened parking lots, thus reducing the load on the most demanded ones.
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>

	<div class="self-center divider w-96 p-12" />
	<div>
		<div class="pb-10">
			<h2 class="font-extrabold text-4xl">Are you a parking operator?</h2>
		</div>
		<a href={mailFromOperator} class="btn btn-primary">Contact us</a>
	</div>

	<div class="pb-12" />
</div>
