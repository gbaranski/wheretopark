<script lang="ts">
	import '/node_modules/flag-icons/css/flag-icons.min.css';

	export let languageSwitcher: boolean = false;
</script>

<li>
	<span class="flex flex-row content-between">
		<a href="/faq" class="basis-full font-mono text-sm">FAQ</a>
		<svg
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			stroke-width="1.5"
			stroke="currentColor"
			class="w-8 h-8"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				d="M9.879 7.519c1.171-1.025 3.071-1.025 4.242 0 1.172 1.025 1.172 2.687 0 3.712-.203.179-.43.326-.67.442-.745.361-1.45.999-1.45 1.827v.75M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9 5.25h.008v.008H12v-.008z"
			/>
		</svg>
	</span>
</li>
<li>
	<span class="flex flex-row content-between">
		<a href="/about" class="basis-full font-mono text-sm">About</a>
		<svg
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			stroke-width="1.5"
			stroke="currentColor"
			class="w-8 h-8"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				d="M4.26 10.147a60.436 60.436 0 00-.491 6.347A48.627 48.627 0 0112 20.904a48.627 48.627 0 018.232-4.41 60.46 60.46 0 00-.491-6.347m-15.482 0a50.57 50.57 0 00-2.658-.813A59.905 59.905 0 0112 3.493a59.902 59.902 0 0110.399 5.84c-.896.248-1.783.52-2.658.814m-15.482 0A50.697 50.697 0 0112 13.489a50.702 50.702 0 017.74-3.342M6.75 15a.75.75 0 100-1.5.75.75 0 000 1.5zm0 0v-3.675A55.378 55.378 0 0112 8.443m-7.007 11.55A5.981 5.981 0 006.75 15.75v-1.5"
			/>
		</svg>
	</span>
</li>
{#if languageSwitcher}
	<li>
		<span class="font-mono text-sm flex flex-row content-between">
			<span class="basis-full"> Switch to Polish </span>
			<span class="fi fi-pl right-[2px]"></span>
		</span>
	</li>
{/if}
