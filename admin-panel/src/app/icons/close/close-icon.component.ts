import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-close-icon',
  template: `
    <svg
      [ngClass]="className"
      viewBox="0 0 512 512"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fill="none"
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="32"
        d="M368 368L144 144m224 0L144 368"
      ></path>
    </svg>
  `,
})
export class CloseIconComponent {
  @Input() className = '';
}
