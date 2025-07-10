<script lang="ts">
	import * as z from 'zod';
	import * as Card from '$lib/components/ui/card/index';
	import { Badge } from '$lib/components/ui/badge/index';
	import { EventRead } from '$lib/schemas/event';

	type EventRead = z.infer<typeof EventRead>;

	let {
		event,
		showFullDescription = false
	}: {
		event: EventRead;
		showFullDescription?: boolean;
	} = $props();

	const today = new Date();
</script>

<Card.Root class="hover:bg-gray-100">
	<Card.Header class="text-left">
		<Card.Title>
			{event.title}
		</Card.Title>
		{#if event.description}
			<Card.Description>
				{!showFullDescription && event.description.length > 80
					? event.description.slice(0, 80) + '...'
					: event.description}
			</Card.Description>
		{/if}
	</Card.Header>
	<Card.Content class="flex gap-x-2 align-middle">
		{#if new Date(event.endDate) < today}
			<Badge variant="secondary">Прошло</Badge>
		{:else if new Date(event.beginDate) <= today && today <= new Date(event.endDate)}
			<Badge>В процессе</Badge>
		{/if}
		{new Date(event.beginDate).toLocaleDateString('ru-RU')} —
		{new Date(event.endDate).toLocaleDateString('ru-RU')}
	</Card.Content>
</Card.Root>
