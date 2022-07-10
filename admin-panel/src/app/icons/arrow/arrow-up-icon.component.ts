import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-arrow-up-icon',
  template: `
    <svg
      viewBox="0 0 24 24"
      [ngClass]="className"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path d="M7 14l5-5 5 5H7z"></path>
    </svg>
  `,
})
export class ArrowUpIconComponent {
  @Input() className = '';
}
