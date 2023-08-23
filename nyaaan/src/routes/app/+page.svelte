<script lang="ts">
	export let data;
</script>

<div class="flex flex-row min-w-screen min-h-screen bg-stone-950">
	<div class="w-16">
		{#each data.servers as server}
			<button class="mx-2 my-2">
				<img
					src="https://avatars.githubusercontent.com/u/46141275?v=4"
					alt="nekonako"
					class="relative inline-block rounded-full"
				/>
			</button>
		{/each}
	</div>
	<div class="w-1/5 border border-stone-900 flex flex-col justify-left">
		{#each data.channels as channel}
			<button class="text-white">
				{channel.name}
			</button>
		{/each}
	</div>
	<div>2</div>
</div>
