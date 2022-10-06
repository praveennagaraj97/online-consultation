import { NgModule } from '@angular/core';
import { StoreModule } from '@ngrx/store';
import { StoreDevtoolsModule } from '@ngrx/store-devtools';
import { environment } from 'src/environments/environment';
import { authReducer } from './auth/auth.reducer';

export interface StoreState {
  auth: { isLogged: boolean };
}

@NgModule({
  imports: [
    StoreModule.forRoot({ auth: authReducer }),
    !environment.production ? StoreDevtoolsModule.instrument() : [],
  ],
})
export class AppStoreModule {}
