import { OverlayModule } from '@angular/cdk/overlay';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ComponentsModule } from './components/component.module';
import { AuthorizedGuard } from './guards/is-authorized.guard';
import { NotAuthorizedGuard } from './guards/not-authorized.guard';
import { IconsModule } from './icons/icons.module';
import { APiInterceptor } from './interceptors/api.interceptor';
import { LayoutModule } from './layouts/layout.module';
import { AppStoreModule } from './store/store.module';
import { ForgotPasswordViewComponent } from './views/auth/forgot-password/forgot-password-view.component';

import { LoginViewComponent } from './views/auth/login/login-view.component';
import { ResetPasswordViewComponent } from './views/auth/reset-password/reset-password-view.component';
import { DashboardViewComponent } from './views/dashboard/dashboard-view.component';
import { PageNotFoundViewComponent } from './views/not-found/not-found-view.component';

@NgModule({
  declarations: [
    AppComponent,
    PageNotFoundViewComponent,
    LoginViewComponent,
    ForgotPasswordViewComponent,
    ResetPasswordViewComponent,
    DashboardViewComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    BrowserAnimationsModule,
    ComponentsModule,
    IconsModule,
    FormsModule,
    AppStoreModule,
    LayoutModule,
    OverlayModule,
  ],
  bootstrap: [AppComponent],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: APiInterceptor,
      multi: true,
    },
    AuthorizedGuard,
    NotAuthorizedGuard,
  ],
})
export class AppModule {}
