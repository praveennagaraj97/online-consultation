import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { LockIconComponent } from './lock/lock-icon.component';
import { MailIconComponent } from './mail/mail-icon.component';
import { SpinnerIconComponent } from './spinner/spinner-icon.component';

@NgModule({
  imports: [CommonModule],
  declarations: [MailIconComponent, LockIconComponent, SpinnerIconComponent],
  exports: [MailIconComponent, LockIconComponent, SpinnerIconComponent],
})
export class IconsModule {}
