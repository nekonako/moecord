<script lang="ts">
	import { goto } from '$app/navigation';
	import { redirect } from '@sveltejs/kit';
	const baseURL = '/api/login/oauth/authorization/';
	const providers = ['google', 'github', 'discord'];
	function handlleOauth(provider: string) {
		goto(baseURL + provider);
	}
</script>

<div class="flex items-center justify-center min-w-screen min-h-screen">
	<div class="flex flex-col items-center">
		<p>Login with</p>
		<div>
			{#each providers as provider}
				<button
					class="px-2 py-1 bg-gray-200 m-2 rounded-lg"
					on:click={() => handlleOauth(provider)}
				>
					{provider}
				</button>
			{/each}
		</div>
	</div>
</div>
