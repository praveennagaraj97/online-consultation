import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  template: `
    <router-outlet></router-outlet>
    <app-theme-provider-component></app-theme-provider-component>
  `,
})
export class AppComponent {}
