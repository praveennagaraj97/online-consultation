import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-active-status-icon',
  template: ` <div [ngClass]="className">Active</div>`,
})
export class ActiveStatusIconComponent {
  @Input() className =
    'bg-green-500 p-1 rounded-xl shadow-md text-center text-gray-50';
}
