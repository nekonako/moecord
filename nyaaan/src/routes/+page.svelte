<script lang="ts">
	import { goto } from '$app/navigation';
</script>

<button on:click={() => goto('/oauth')} class="p-4 ring-slate-200 ring-2"> Login </button>
