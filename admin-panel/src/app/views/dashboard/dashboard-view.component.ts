import { Component } from '@angular/core';
import { Store } from '@ngrx/store';
import { AuthState } from 'src/app/store/auth/auth.types';

@Component({
  selector: 'app-dashboard-view',
  template: ` <h1>{{ isLogged }}</h1>`,
})
export class DashboardViewComponent {
  isLogged = false;

  constructor(private store: Store<{ auth: AuthState }>) {}

  ngOnInit() {
    this.store.select('auth').subscribe({
      next: (val) => {
        this.isLogged = val.isLogged;
      },
    });
  }
}
