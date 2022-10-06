export function addDaysToDate(days: number, dateCursor?: Date) {
  const date = dateCursor || new Date();

  date.setDate(date.getDate() + days);

  return date;
}

export function addMinutesToDate(mins: number, dateCursor?: Date) {
  const date = dateCursor || new Date();

  date.setMinutes(date.getMinutes() + mins);

  return date;
}

export function subtractMinutes(
  minutes: number,
  dateCursor: Date = new Date()
) {
  dateCursor.setMinutes(dateCursor.getMinutes() - minutes);

  return dateCursor;
}
