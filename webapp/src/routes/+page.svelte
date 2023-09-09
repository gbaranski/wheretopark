<script lang="ts">
	import { getRandomBetween, mailFromOperator } from '$/lib/utils';
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

	const cities = ['gdansk', 'gdynia', 'glasgow', 'klodzko', 'poznan', 'sopot', 'warsaw'];

	const awards = [
		{
			name: 'explory',
			description: 'Finalist of a high-tech contest for innovative ideas',
			link: 'https://www.explory.pl/2023/naukowcy/32-grzegorz-baranski/'
		},
		{
			name: 'young-innovator',
			description: 'III Place in Young Innovator 2022/2023 for HS Students',
			link: 'https://not.org.pl/olimpiady-i-konkursy/xvi-edycja-konkursu-mlody-innowator-20222023?department=centrala'
		}
	];
	
	const news = [
		{
			name: 'eska.pl',
			link: 'https://www.eska.pl/trojmiasto/licealista-z-gdanska-tworzy-aplikacje-ktora-ulatwi-znalezienie-miejsca-na-parkingu-aa-MCvr-56JZ-Gnus.html',
		},
		{
			name: 'otwartedane.gdynia.pl',
			link: 'https://otwartedane.gdynia.pl/where-to-park-aplikacja-do-szukania-wolnych-miejsc-parkingowych/',
		},
		{
			name: 'mlodagdynia.pl',
			link: 'https://mlodagdynia.pl/pl/19_wiadomosci-z-regionu/53682_wolne-miejsca-na-parkingu-pomoze-aplikacja-trojmiejskiego-licealisty.html',
		},
		{
			name: 'blog.citydata.pl',
			link: 'https://blog.citydata.pl/where-to-park-aplikacja-do-szukania-wolnych-miejsc-parkingowych/',
		}
	];
</script>

<div class="flex flex-col text-center w-full pt-24 p-12">
	<h1 class="font-extrabold text-5xl lg:text-7xl">
		Where <span class="text-primary">To</span> Park
	</h1>
	<div class="self-center divider w-32" />
	<h2 class="font-regular text-lg lg:text-xl">
		With the help of our AI we'll find you an available parking spot nearby!
	</h2>
	<br />
	<a
		href="https://apps.apple.com/us/app/where-to-park/id6444453582?itsct=apps_box_badge&amp;itscg=30200"
		class="block overflow-hidden self-center"
		><img
			src="https://tools.applemediaservices.com/api/badges/download-on-the-app-store/white/en-us?size=250x83&amp;releaseDate=1668988800&h=a0d7d9ddd291f5ab79945e444e83a9f9"
			alt="Download on the App Store"
			class="w-48"
		/></a
	>
	<p class="pt-2">
		or
		<a class="link link-info" href="/app">open app in the browser</a>
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
		<img class="w-96 lg:w-2/3 rounded-2xl" src="/assets/preview.webp" alt="preview of the app" />
	</div>

	<div class="self-center divider w-96 p-12" />
	<div class="flex justify-evenly items-center flex-col-reverse lg:flex-row gap-8">
		<div class="">
			<h2 class="font-extrabold text-4xl">Powered by AI &nbsp;🤖</h2>
			Our system automatically marks<span class="text-green-600">green</span> free spaces, and
			<span class="text-red-600">red</span> occupied.
		</div>
		<img
			class="w-96 lg:w-2/3 rounded-lg"
			src="/assets/parking-basen-klodzko.webp"
			alt="animation of our ai"
		/>
	</div>

	<div class="self-center divider w-96 p-12" />
	<h2 class="font-extrabold text-4xl">We have parkings lots in</h2>
	<div class="pt-10 flex flex-row flex-wrap gap-16 justify-center">
		{#each cities as city}
			<img src={`/assets/city/${city}.webp`} alt="{city} logo" class="w-32 object-scale-down" />
		{/each}
	</div>

	<div class="self-center divider w-96 p-12" />
	<h2 class="font-extrabold text-4xl">Awards</h2>
	<div class="pt-10 flex flex-row flex-wrap gap-16 justify-center">
		{#each awards as award}
			<a class="w-32" href="{award.link}" target="_blank">
				<img
					src={`/assets/award/${award.name}.webp`}
					alt="{award.name} logo"
					class="h-32 object-scale-down pb-2"
				/>
				<span class="link link-secondary">{award.description}</span>
			</a>
		{/each}
	</div>

	<div class="self-center divider w-96 p-12" />
	<h2 class="font-extrabold text-4xl">Talks about us</h2>
	<div class="pt-10 flex flex-row flex-wrap gap-16 justify-center">
		{#each news as entry}
			<a class="w-32" href="{entry.link}" target="_blank">
				<img
					src={`/assets/news/${entry.name}.webp`}
					alt="{entry.name} logo"
					class="h-32 object-scale-down"
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
