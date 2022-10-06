import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  template: `
    <router-outlet></router-outlet>
    <app-theme-provider></app-theme-provider>
  `,
})
export class AppComponent {
  constructor() {}
}
