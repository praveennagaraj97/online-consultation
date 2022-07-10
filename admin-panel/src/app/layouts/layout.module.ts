import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { IconsModule } from '../components/icons/icons.module';
import { HeaderComponent } from './protected/header/header.component';
import { ProtectedLayoutComponent } from './protected/layout.component';
import { NavBarComponent } from './protected/nav-bar/nav-bar.component';
import { NavItemComponent } from './protected/nav-bar/nav-item/nav-item.component';

@NgModule({
  declarations: [
    ProtectedLayoutComponent,
    HeaderComponent,
    NavBarComponent,
    NavItemComponent,
  ],
  exports: [ProtectedLayoutComponent],
  imports: [CommonModule, RouterModule, IconsModule],
})
export class LayoutModule {}
