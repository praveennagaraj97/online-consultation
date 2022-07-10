import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-dots-icon',
  template: `
    <svg
      viewBox="0 0 16 16"
      [ngClass]="className"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        d="M3 9.5a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3zm5 0a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3zm5 0a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3z"
      ></path>
    </svg>
  `,
})
export class DotsIconComponent {
  @Input() className = '';
}
