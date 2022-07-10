import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { ComponentsModule } from '../components/component.module';
import { DirectiveModule } from '../directives/directive.module';
import { IconsModule } from '../icons/icons.module';
import { PipesModule } from '../pipes/pipes.module';
import { HeaderComponent } from './protected/header/header.component';
import { UserOptionsDropdownComponent } from './protected/header/user/options-dropdown/dropdown.component';
import { UserOptionsComponent } from './protected/header/user/user-options.component';
import { ProtectedLayoutComponent } from './protected/layout.component';
import { NavBarComponent } from './protected/nav-bar/nav-bar.component';
import { NavItemComponent } from './protected/nav-bar/nav-item/nav-item.component';

@NgModule({
  declarations: [
    ProtectedLayoutComponent,
    HeaderComponent,
    NavBarComponent,
    NavItemComponent,
    UserOptionsComponent,
    UserOptionsDropdownComponent,
  ],
  exports: [ProtectedLayoutComponent],
  imports: [
    CommonModule,
    RouterModule,
    IconsModule,
    ComponentsModule,
    PipesModule,
    DirectiveModule,
  ],
})
export class LayoutModule {}
