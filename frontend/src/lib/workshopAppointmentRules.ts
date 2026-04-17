/** Workshop rules: Mon–Fri, morning 9:30–12:30, afternoon 14:00–17:30 (local), 30-minute slots, max 8 non-cancelled per calendar day (server). */

export const WORKSHOP_MAX_APPOINTMENTS_PER_DAY = 8;

const morningStart = 9 * 60 + 30;
const morningEnd = 12 * 60 + 30;
const afternoonStart = 14 * 60;
const afternoonEnd = 17 * 60 + 30;

export function pad2(n: number): string {
  return n < 10 ? `0${n}` : String(n);
}

function minutesToHHmm(totalMinutes: number): string {
  const h = Math.floor(totalMinutes / 60);
  const m = totalMinutes % 60;
  return `${pad2(h)}:${pad2(m)}`;
}

/** Every 30 minutes from start through end (inclusive). */
function halfHourSlotsInclusive(startMin: number, endMin: number): string[] {
  const out: string[] = [];
  for (let m = startMin; m <= endMin; m += 30) {
    out.push(minutesToHHmm(m));
  }
  return out;
}

/** All allowed clock times (HH:mm, 24h), excluding weekday / “today” filtering. */
export const ALL_WORKSHOP_SLOT_TIMES_HHMM: readonly string[] = [
  ...halfHourSlotsInclusive(morningStart, morningEnd),
  ...halfHourSlotsInclusive(afternoonStart, afternoonEnd),
];

function clockMinutes(d: Date): number {
  return d.getHours() * 60 + d.getMinutes();
}

/** True if local wall time is inside allowed bands (inclusive). */
export function isWithinWorkshopHours(d: Date): boolean {
  const m = clockMinutes(d);
  return (m >= morningStart && m <= morningEnd) || (m >= afternoonStart && m <= afternoonEnd);
}

export function workshopHoursErrorMessage(): string {
  return 'O horário tem de estar entre as 9:30 e as 12:30 ou entre as 14:00 e as 17:30.';
}

/** Parse YYYY-MM-DD as local calendar date (no UTC shift). */
export function parseLocalDateOnlyYmd(ymd: string): Date | null {
  const m = /^(\d{4})-(\d{2})-(\d{2})$/.exec(ymd.trim());
  if (!m) return null;
  const y = Number(m[1]);
  const mo = Number(m[2]) - 1;
  const d = Number(m[3]);
  const dt = new Date(y, mo, d, 12, 0, 0, 0);
  if (dt.getFullYear() !== y || dt.getMonth() !== mo || dt.getDate() !== d) return null;
  return dt;
}

/** Monday (1) through Friday (5) in local time. */
export function isWeekdayLocalYmd(ymd: string): boolean {
  const d = parseLocalDateOnlyYmd(ymd);
  if (!d) return false;
  const day = d.getDay();
  return day >= 1 && day <= 5;
}

export function weekdayErrorMessage(): string {
  return 'Só é possível marcar de segunda a sexta-feira.';
}

/** Local YYYY-MM-DD for today. */
export function todayLocalYmd(): string {
  const n = new Date();
  return `${n.getFullYear()}-${pad2(n.getMonth() + 1)}-${pad2(n.getDate())}`;
}

export function localYmdFromDate(d: Date): string {
  return `${d.getFullYear()}-${pad2(d.getMonth() + 1)}-${pad2(d.getDate())}`;
}

/** Calendar day YYYY-MM-DD for an ISO-like string, in local time. */
export function localYmdFromIso(iso: string): string {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return '';
  return localYmdFromDate(d);
}

/**
 * Picks the nearest workshop slot (by clock distance) to the local time of `iso`.
 * Returns '' if the instant is too far from any allowed slot (e.g. wrong day context).
 */
export function workshopSlotFromIso(iso: string): string {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return '';
  const raw = d.getHours() * 60 + d.getMinutes();
  let best = '';
  let bestDist = 99999;
  for (const s of ALL_WORKSHOP_SLOT_TIMES_HHMM) {
    const [h, mi] = s.split(':').map(Number);
    const sm = h * 60 + mi;
    const dist = Math.abs(sm - raw);
    if (dist < bestDist) {
      bestDist = dist;
      best = s;
    }
  }
  return bestDist <= 24 ? best : '';
}

/** Combine local date + HH:mm into a string suitable for APIs (no timezone suffix). */
export function combineLocalDateAndSlot(ymd: string, hhmm: string): string {
  return `${ymd}T${hhmm}`;
}

/**
 * Slots still bookable for the given calendar day: workshop list minus times not strictly after `now`.
 * `bufferMinutes` adds a small margin after “now” (e.g. 30).
 */
export function getBookableSlotTimesForDay(
  ymd: string,
  now: Date = new Date(),
  bufferMinutes = 30,
): string[] {
  const day = parseLocalDateOnlyYmd(ymd);
  if (!day || !isWeekdayLocalYmd(ymd)) return [];

  const cutoff = new Date(now.getTime() + bufferMinutes * 60 * 1000);
  const isSameLocalDay =
    cutoff.getFullYear() === day.getFullYear() &&
    cutoff.getMonth() === day.getMonth() &&
    cutoff.getDate() === day.getDate();

  if (!isSameLocalDay) {
    return [...ALL_WORKSHOP_SLOT_TIMES_HHMM];
  }

  const limitM = cutoff.getHours() * 60 + cutoff.getMinutes();
  return ALL_WORKSHOP_SLOT_TIMES_HHMM.filter((hhmm) => {
    const [h, mi] = hhmm.split(':').map(Number);
    const slotM = h * 60 + mi;
    return slotM > limitM;
  });
}

/** Label for UI: 24h clock + “horas” (e.g. 14:00 horas). */
export function formatSlotLabel24h(hhmm: string): string {
  return `${hhmm} horas`;
}

/** Summary line: weekday date + time in 24h. */
export function formatAppointmentSummaryLocal(ymd: string, hhmm: string): string {
  const d = parseLocalDateOnlyYmd(ymd);
  if (!d) return `${ymd} ${hhmm}`;
  const datePart = new Intl.DateTimeFormat('pt-PT', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  }).format(d);
  return `${datePart}, ${hhmm} horas`;
}
