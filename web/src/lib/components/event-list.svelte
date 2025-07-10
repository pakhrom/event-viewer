<script lang="ts">
	import * as z from 'zod';
	import { MediaQuery } from 'svelte/reactivity';
	import * as Tabs from '$lib/components/ui/tabs/index';
	import * as Dialog from '$lib/components/ui/dialog/index';
	import * as Drawer from '$lib/components/ui/drawer/index';
	import { Badge } from '$lib/components/ui/badge';
	import { Skeleton } from '$lib/components/ui/skeleton/index';
	import { buttonVariants } from '$lib/components/ui/button/index';
	import { EventsRead } from '$lib/schemas/event';
	import Event from './event.svelte';

	type EventsRead = z.infer<typeof EventsRead>;

	let {
		events,
		showFilter = false,
		showEditButton = false
	}: {
		events: EventsRead;
		showFilter?: boolean;
		showEditButton?: boolean;
	} = $props();

	let infoOpen: boolean = $state(false);
	let selectedEventIndex: number = $state(0);
	const isDesktop = new MediaQuery('(min-width: 768px)');

	let today = new Date();
</script>

{#if showFilter}
	<Tabs.Root value="future">
		<Tabs.List>
			<Tabs.Trigger
				value="past"
				onclick={() => {
					console.log('test');
				}}
			>
				Прошедшие
			</Tabs.Trigger>
			<Tabs.Trigger value="present">Текущие</Tabs.Trigger>
			<Tabs.Trigger value="future">Будущие</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content value="past">
			<ol class="flex flex-col gap-4">
				<Skeleton class="h-[8em] max-w-prose rounded-2xl" />
				<Skeleton class="h-[8em] max-w-prose rounded-2xl" />
				<Skeleton class="h-[8em] max-w-prose rounded-2xl" />
			</ol>
		</Tabs.Content>
		<Tabs.Content value="present">
			<ol class="flex flex-col gap-4">
				<Skeleton class="h-[8em] max-w-prose rounded-2xl" />
				<Skeleton class="h-[8em] max-w-prose rounded-2xl" />
				<Skeleton class="h-[8em] max-w-prose rounded-2xl" />
			</ol>
		</Tabs.Content>
		<Tabs.Content value="future">
			<ol class="flex flex-col gap-4">
				<Skeleton class="h-[8em] max-w-prose rounded-2xl" />
				<Skeleton class="h-[8em] max-w-prose rounded-2xl" />
				<Skeleton class="h-[8em] max-w-prose rounded-2xl" />
			</ol>
		</Tabs.Content>
	</Tabs.Root>
{:else}
	<ol class="flex flex-col gap-4">
		{#each events as event, eventIndex (event.id)}
			<button
				class="hover:cursor-pointer"
				type="button"
				onclick={() => {
					infoOpen = true;
					selectedEventIndex = eventIndex;
				}}
			>
				<Event {event}></Event>
			</button>
		{/each}
	</ol>
{/if}

{#if isDesktop.current}
	<Dialog.Root bind:open={infoOpen}>
		{#if events[selectedEventIndex]}
			<Dialog.Content class="sm:max-w-2/3">
				<Dialog.Header>
					<Dialog.Title class="flex gap-2 align-middle">
						{#if new Date(events[selectedEventIndex].endDate) < today}
							<Badge variant="secondary">Прошло</Badge>
						{:else if new Date(events[selectedEventIndex].beginDate) <= today && today <= new Date(events[selectedEventIndex].endDate)}
							<Badge>В процессе</Badge>
						{/if}
						{events[selectedEventIndex].title}
					</Dialog.Title>
					<Dialog.Description>
						{events[selectedEventIndex].description}
					</Dialog.Description>
				</Dialog.Header>
				<Dialog.Close class={buttonVariants({ variant: 'outline' })}>Редактировать</Dialog.Close>
			</Dialog.Content>
		{/if}
	</Dialog.Root>
{:else}
	<Drawer.Root bind:open={infoOpen}>
		{#if events[selectedEventIndex]}
			<Drawer.Content>
				<Drawer.Header class="text-left">
					<Drawer.Title>
						{#if new Date(events[selectedEventIndex].endDate) < today}
							<Badge variant="secondary">Прошло</Badge>
						{:else if new Date(events[selectedEventIndex].beginDate) <= today && today <= new Date(events[selectedEventIndex].endDate)}
							<Badge>В процессе</Badge>
						{/if}
						{events[selectedEventIndex].title}
					</Drawer.Title>
					<Drawer.Description>
						{events[selectedEventIndex].description}
					</Drawer.Description>
				</Drawer.Header>
				{#if showEditButton}
					<Drawer.Footer class="pt-2">
						<Drawer.Close class={buttonVariants({ variant: 'outline' })}>
							Редактировать
						</Drawer.Close>
					</Drawer.Footer>
				{/if}
			</Drawer.Content>
		{/if}
	</Drawer.Root>
{/if}
