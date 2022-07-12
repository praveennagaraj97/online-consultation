import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-next-icon',
  template: `
    <svg
      viewBox="0 0 24 24"
      [ngClass]="className"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fill="none"
        class="dark:stroke-gray-50 stroke-gray-900"
        stroke-width="2"
        d="M6,12.4 L18,12.4 M12.6,7 L18,12.4 L12.6,17.8"
      ></path>
    </svg>
  `,
})
export class NextIconComponent {
  @Input() className = '';
}
