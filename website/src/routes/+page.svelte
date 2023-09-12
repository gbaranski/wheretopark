<script lang="ts">
	import { awards } from '$lib/assets/award';
	import { cities } from '$lib/assets/city';
	import { news } from '$lib/assets/news';
	import appPreview from '$lib/assets/preview.webp';
	import exampleParkingLot from '$lib/assets/parking-basen-klodzko.webp';
	import { getRandomBetween, mailFromOperator } from '$lib/utils';
	import { inview } from 'svelte-inview';

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
</script>

<div class="flex flex-col text-center w-full pt-24 p-12">
	<div class="flex items-center justify-center">
		<div class="w-max">
			<h1
				class="animate-typing overflow-hidden whitespace-nowrap border-r-4 border-r-secondary pr-5 text-2xl lg:text-5xl font-extrabold lg:leading-[4rem]"
			>
				Your parking assistant &nbsp;
			</h1>
		</div>
	</div>
	<h2 class="font-regular text-lg lg:text-xl pt-2">
		With the help of our AI we'll find you <b>an available</b> parking spot nearby!
	</h2>
	<br />
	<div>
		<a
			href="https://apps.apple.com/us/app/where-to-park/id6444453582?itsct=apps_box_badge&amp;itscg=30200"
			style="display: inline-block; overflow: hidden; border-radius: 13px; width: 250px; height: 83px;"
			><img
				src="https://tools.applemediaservices.com/api/badges/download-on-the-app-store/white/en-us?size=250x83&amp;releaseDate=1668988800"
				alt="Download on the App Store"
				style="border-radius: 13px; width: 250px; height: 83px;"
			/></a
		>
	</div>
	<p class="pt-2">
		or
		<a class="link link-info" href="https://web.wheretopark.app">open app in the browser</a>
	</p>

	<div class="self-center divider w-96 p-12" />
	<div class="flex justify-evenly items-center flex-col lg:flex-row">
		<h2 class="font-extrabold text-3xl">What problem are we solving?</h2>
		<div class="divider self-center w-8 lg:hidden" />
		<div>
			<h2 class="font-extrabold text-5xl text-info">
				<span use:inview on:inview_enter={onCounterVisible}>{counter}</span> minutes
			</h2>
			<p>
				The average amount of <span class="font-bold">time wasted</span> each time you want to park.
			</p>
		</div>
	</div>

	<div class="self-center divider w-96 p-12" />
	<div class="flex justify-evenly items-center flex-col lg:flex-row-reverse gap-8">
		<div class="">
			<h2 class="font-extrabold text-4xl">Our solution</h2>
			<h3 class="">Easy access to information about available parking lots nearby</h3>
		</div>
		<img class="w-96 lg:w-2/3 rounded-2xl" src={appPreview} alt="preview of the app" width="384px" height="207px" />
	</div>

	<div class="self-center divider w-96 p-12" />
	<div class="flex justify-evenly items-center flex-col-reverse lg:flex-row gap-8">
		<div class="">
			<h2 class="font-extrabold text -4xl">Powered by AI &nbsp;🤖</h2>
			Our system automatically marks&nbsp;<span class="text-green-600">green</span> free spaces, and
			<span class="text-red-600">red</span> occupied.
		</div>
		<img class="w-96 lg:w-2/3 rounded-lg" src={exampleParkingLot} alt="animation of our ai" width="384px" height="216px"/>
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
			<h2 class="font-extrabold text-4xl">Are you a parking operator?</h2>
		</div>
		<a href={mailFromOperator} class="btn btn-primary">Contact us</a>
	</div>

	<div class="pb-12" />
</div>
