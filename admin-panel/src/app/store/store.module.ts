import { NgModule } from '@angular/core';
import { StoreModule } from '@ngrx/store';
import { StoreDevtoolsModule } from '@ngrx/store-devtools';
import { environment } from 'src/environments/environment';
import { authReducer } from './auth/auth.reducer';
import { AuthState } from './auth/auth.types';

export interface StoreState {
  auth: AuthState;
}

@NgModule({
  imports: [
    StoreModule.forRoot({ auth: authReducer }),
    !environment.production ? StoreDevtoolsModule.instrument() : [],
  ],
})
export class AppStoreModule {}
