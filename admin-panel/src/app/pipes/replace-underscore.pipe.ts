import { Pipe, PipeTransform } from '@angular/core';

@Pipe({ name: 'replaceUnderscore' })
export class ReplaceUnderScorePipe implements PipeTransform {
  transform(value: string) {
    return value.replace(/\_/g, ' ');
  }
}
