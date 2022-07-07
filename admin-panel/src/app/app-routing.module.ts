import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ForgotPasswordViewComponent } from './views/auth/forgot-password/forgot-password.component';
import { LoginViewComponent } from './views/auth/login/login-view.component';
import { ResetPasswordViewComponent } from './views/auth/reset-password/reset-password-view.component';
import { PageNotFoundViewComponent } from './views/not-found/not-found-view.component';

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
      {
        path: 'reset-password',
        component: ResetPasswordViewComponent,
      },
    ],
  },
];

const protectedRoues: Routes = [];

@NgModule({
  imports: [
    RouterModule.forRoot([
      ...publicRoues,
      ...protectedRoues,
      { path: '**', component: PageNotFoundViewComponent },
    ]),
  ],
  exports: [RouterModule],
})
export class AppRoutingModule {}
