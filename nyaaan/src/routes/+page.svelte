<script lang="ts">
	let button = 0;
	function incrmenet() {
		button += 1;
	}
</script>

<h1 class="underline">Hello world!</h1>
<button on:click={incrmenet} class="p-4 ring-slate-200 ring-2">{button}</button>
