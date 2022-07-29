import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-inactive-status-icon',
  template: ` <div [ngClass]="className">In Active</div>`,
})
export class InActiveStatusIconComponent {
  @Input() className =
    'bg-red-500 p-1 rounded-xl shadow-md text-center text-gray-50';
}
