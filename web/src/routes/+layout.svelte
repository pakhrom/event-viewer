<script lang="ts">
	import '../app.css';

	let { children } = $props();
</script>

<div class="mx-auto mt-12 max-w-prose">
	<div class="mx-2">
		{@render children()}
	</div>
</div>

<style>
</style>
