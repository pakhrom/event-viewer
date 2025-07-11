<script lang="ts">
	import type z from 'zod';
	import { EventsRead } from '$lib/schemas/event';
	import SuperDebug from 'sveltekit-superforms';
	import EventList from '$lib/components/event-list.svelte';

	type EventsRead = z.infer<typeof EventsRead>;

	let eventsRead: string =
		'[{"id":"1","title":"Спринт по разработке веб сервиса","description":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip  fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt iLorem ipsum dolor sit amet, consectetur adipiscing elit,  ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation  Duis aute irure dolorn culpa qui officia deserunt mollit anim id est laborum.","beginDate":"2025-07-10","endDate":"2025-07-12"},{"id":"12","title":"Спринт по разработке веб сервиса","description":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.","beginDate":"2025-07-10","endDate":"2025-07-12"},{"id":"2","title":"Вечеринка в честь Нового 2025 Года","description":"Записать можно по ссылке: https://youtu.be/dQw4w9WgXcQ?si=q7D14EFC3RvroPaK","beginDate":"2024-12-29","endDate":"2025-01-02"},{"id":"3","title":"Снова в шарагу","description":"школа никого не пощадит","beginDate":"2025-09-01","endDate":"2025-09-03"}]';
	const eventsParsed = $state(EventsRead.safeParse(JSON.parse(eventsRead)));
	let events: EventsRead | undefined = $state();
	if (!eventsParsed.success) {
		console.error(eventsParsed.error);
	} else {
		events = eventsParsed.data;
	}
</script>

<main>
	<h1 class="text-3xl font-bold text-balance">Учебные мероприятия</h1>

	{#if events}
		<EventList {events} showEditButton></EventList>
	{:else}
		Не удалось загрузить список мероприятий. <br />
		<SuperDebug data={eventsParsed.error?.issues}></SuperDebug>
	{/if}
</main>

<style>
	h1 {
		margin-bottom: 0.3em;
	}
</style>
