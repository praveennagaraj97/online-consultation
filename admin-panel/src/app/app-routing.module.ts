import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ForgotPasswordViewComponent } from './views/auth/forgot-password/forgot-password.component';
import { LoginViewComponent } from './views/auth/login/login-view.component';

const publicRoues: Routes = [
  {
    path: 'auth',
    children: [
      {
        path: '',
        redirectTo: 'login',
        pathMatch: 'full',
      },
      {
        path: 'login',
        component: LoginViewComponent,
      },
      {
        path: 'forgot-password',
        component: ForgotPasswordViewComponent,
      },
    ],
  },
];

const protectedRoues: Routes = [];

@NgModule({
  imports: [RouterModule.forRoot([...publicRoues, ...protectedRoues])],
  exports: [RouterModule],
})
export class AppRoutingModule {}
