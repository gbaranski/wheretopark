<script lang="ts">
	export let languageSwitcher: boolean = false;
</script>

<ul class="menu menu-horizontal px-1">
	<li>
		<a href="/faq" class="basis-full text-sm"
			>FAQ
		</a>
	</li>
	<li>
		<a href="/about" class="basis-full text-sm"
			>About
		</a>
	</li>
	{#if languageSwitcher}
		<li>
			<span class="font-mono text-sm flex flex-row content-between">
				<span class="basis-full"> Switch to Polish </span>
				<span class="fi fi-pl right-[2px]" />
			</span>
		</li>
	{/if}
</ul>