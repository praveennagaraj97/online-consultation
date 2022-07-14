import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-search-icon',
  template: `
    <svg
      viewBox="0 0 512 512"
      [ngClass]="className"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fill="none"
        stroke-miterlimit="10"
        stroke-width="32"
        d="M221.09 64a157.09 157.09 0 10157.09 157.09A157.1 157.1 0 00221.09 64z"
      ></path>
      <path
        fill="none"
        stroke-linecap="round"
        stroke-miterlimit="10"
        stroke-width="32"
        d="M338.29 338.29L448 448"
      ></path>
    </svg>
  `,
})
export class SearchIconComponent {
  @Input() className = '';
}
