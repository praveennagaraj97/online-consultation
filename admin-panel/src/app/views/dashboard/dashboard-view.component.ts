import { Component } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { Store } from '@ngrx/store';
import { StoreState } from 'src/app/store/store.module';

@Component({
  selector: 'app-dashboard-view',
  templateUrl: 'dashboard-view.component.html',
})
export class DashboardViewComponent {
  isLogged = false;

  constructor(private store: Store<StoreState>, private titleService: Title) {}

  ngOnInit() {
    this.store.select('auth').subscribe({
      next: (val) => {
        this.isLogged = val.isLogged;
      },
    });
  }
}
