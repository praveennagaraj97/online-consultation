import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-chevron-icon',
  template: `
    <svg
      viewBox="0 0 24 24"
      [ngClass]="className"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        d="M10.707 17.707 16.414 12l-5.707-5.707-1.414 1.414L13.586 12l-4.293 4.293z"
      ></path>
    </svg>
  `,
})
export class ChevronIconComponent {
  @Input() className = '';
}
