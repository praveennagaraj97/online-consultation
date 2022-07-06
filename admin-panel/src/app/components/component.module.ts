import { OverlayModule } from '@angular/cdk/overlay';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { ThemeProviderComponent } from './providers/theme-provider/theme-provider.component';

@NgModule({
  declarations: [ThemeProviderComponent],
  exports: [ThemeProviderComponent],
  imports: [CommonModule, OverlayModule],
})
export class ComponentsModule {}
