import { NgModule } from '@angular/core';
import { PreloadAllModules, RouterModule, Routes } from '@angular/router';
import { AuthorizedGuard } from './guards/is-authorized.guard';
import { NotAuthorizedGuard } from './guards/not-authorized.guard';
import { ProtectedLayoutComponent } from './layouts/protected/layout.component';

import { ForgotPasswordViewComponent } from './views/auth/forgot-password/forgot-password-view.component';
import { LoginViewComponent } from './views/auth/login/login-view.component';
import { ResetPasswordViewComponent } from './views/auth/reset-password/reset-password-view.component';
import { DashboardViewComponent } from './views/dashboard/dashboard-view.component';
import { PageNotFoundViewComponent } from './views/not-found/not-found-view.component';

const publicRoues: Routes = [
  {
    path: 'auth',
    canActivate: [NotAuthorizedGuard],
    children: [
      {
        path: '',
        redirectTo: 'login',
        pathMatch: 'full',
      },
      {
        path: 'login',
        component: LoginViewComponent,
        title: 'Online Consultation | Login',
      },
      {
        path: 'forgot-password',
        component: ForgotPasswordViewComponent,
        title: 'Online Consultation | Forgot Password',
      },
      {
        path: 'reset-password',
        component: ResetPasswordViewComponent,
        title: 'Online Consultation | Reset Password',
      },
    ],
  },
];

const protectedRoues: Routes = [
  {
    path: '',
    component: ProtectedLayoutComponent,
    canActivateChild: [AuthorizedGuard],
    children: [
      { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
      {
        path: 'dashboard',
        component: DashboardViewComponent,
        title: 'Online Consultation | Dashboard',
      },
      {
        path: 'settings',
        component: DashboardViewComponent,
        title: 'Online Consultation | Settings',
      },
      {
        path: 'doctors',
        loadChildren: () =>
          import('./views/doctors/doctors.module').then(
            (m) => m.DoctorsViewModule
          ),
      },
    ],
  },
];

@NgModule({
  imports: [
    RouterModule.forRoot(
      [
        ...publicRoues,
        ...protectedRoues,
        { path: '**', component: PageNotFoundViewComponent },
      ],
      { preloadingStrategy: PreloadAllModules }
    ),
  ],
  exports: [RouterModule],
})
export class AppRoutingModule {}
