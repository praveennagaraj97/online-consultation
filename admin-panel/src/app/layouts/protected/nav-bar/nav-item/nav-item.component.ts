import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-nav-item',
  template: `
    <a
      [routerLink]="path"
      class="flex items-center space-x-3 bg- py-2 px-4 
      rounded-lg text-xs w-full smooth-animate
      hover:bg-sky-700/80 bg-sky-700/40 drop-shadow-2xl shadow-lg border 
      hover:scale-105 border-sky-50/20"
    >
      <ng-content></ng-content>
      <p class="text-sm">{{ title }}</p>
    </a>
  `,
})
export class NavItemComponent {
  @Input() title: string = '';
  @Input() path: string = '';
}
