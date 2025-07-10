import * as z from 'zod';

export const EventWrite = z.object({
  title: z.string().trim().min(1).max(40),
  description: z.string().trim().max(600).optional(),
  beginDate: z.iso.date(),
  endDate: z.iso.date()
})

export const EventRead = EventWrite.extend({
  id: z.string()
})

export const EventsRead = z.array(EventRead)