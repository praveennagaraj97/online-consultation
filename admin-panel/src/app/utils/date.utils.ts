export function addDaysToDate(days: number, dateCursor?: Date) {
  const date = dateCursor || new Date();

  date.setDate(date.getDate() + days);

  return date;
}

export function addMinutesToDate(mins: number, dateCursor?: Date) {
  const date = dateCursor || new Date();

  date.setMinutes(date.getMinutes() + 30);

  return date;
}
