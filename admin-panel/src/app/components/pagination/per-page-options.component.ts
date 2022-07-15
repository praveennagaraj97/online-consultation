import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-per-page',
  template: `
    <div class="flex items-center space-x-1.5">
      <small>Show</small>
      <select
        name="per_page"
        [ngModel]="defaultSelected"
        (ngModelChange)="onChangeListner($event)"
        [value]="defaultSelected"
        class="common-input w-full input-focus input-colors rounded-md p-1 text-sm min-w-fit"
      >
        <option [value]="opt.value" *ngFor="let opt of options">
          {{ opt.title }}
        </option>
      </select>
      <small>Entries</small>
    </div>
  `,
})
export class PerPageOptionsComponent {
  @Input() options: { value: string | number; title: string }[] = [
    { title: '10', value: 10 },
    { title: '20', value: 20 },
    { title: '50', value: 50 },
  ];
  @Input() defaultSelected?: string | number;
  @Output() onChange = new EventEmitter<number>();

  onChangeListner(value: number) {
    this.onChange.emit(value);
  }
}
