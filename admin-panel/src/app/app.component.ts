import { Component } from '@angular/core';
import { Store } from '@ngrx/store';
import { rehydrateAuthState } from './store/auth/auth.actions';
import { StoreState } from './store/store.module';

@Component({
  selector: 'app-root',
  template: `
    <router-outlet></router-outlet>
    <app-theme-provider></app-theme-provider>
  `,
})
export class AppComponent {
  constructor(private store: Store<StoreState>) {
    this.store.dispatch(rehydrateAuthState());
  }
}
