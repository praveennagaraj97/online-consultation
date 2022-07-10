import { Component } from '@angular/core';
import { Store } from '@ngrx/store';
import { rehydrateAuthState } from './store/auth/auth.actions';
import { StoreState } from './store/store.module';

@Component({
  selector: 'app-root',
  template: `
    <router-outlet></router-outlet>
    <app-theme-provider-component></app-theme-provider-component>
  `,
})
export class AppComponent {
  constructor(private store: Store<StoreState>) {
    this.store.dispatch(rehydrateAuthState());
  }
}
