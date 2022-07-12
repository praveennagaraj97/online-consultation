import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-three-dots-icon',
  template: `
    <svg
      viewBox="0 0 24 24"
      [ngClass]="className"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"
      ></path>
    </svg>
  `,
})
export class ThreeDotsIconComponent {
  @Input() className = '';
}
