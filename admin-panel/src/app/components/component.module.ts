import { OverlayModule } from '@angular/cdk/overlay';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { EmailSentComponent } from './email-sent/email-sent.component';
import { ThemeProviderComponent } from './providers/theme-provider/theme-provider.component';

@NgModule({
  imports: [CommonModule, OverlayModule],
  declarations: [ThemeProviderComponent, EmailSentComponent],
  exports: [ThemeProviderComponent, EmailSentComponent],
})
export class ComponentsModule {}
