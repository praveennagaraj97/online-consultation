import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ComponentsModule } from './components/component.module';
import { IconsModule } from './components/icons/icons.module';
import { APiInterceptor } from './interceptors/api.interceptor';
import { LoginViewComponent } from './views/auth/login/login-view.component';

@NgModule({
  declarations: [AppComponent, LoginViewComponent],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    BrowserAnimationsModule,
    ComponentsModule,
    IconsModule,
    FormsModule,
  ],
  bootstrap: [AppComponent],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: APiInterceptor,
      multi: true,
    },
  ],
})
export class AppModule {}
