import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { EmailSentComponent } from './email-sent/email-sent.component';
import { ThemeProviderComponent } from './providers/theme-provider/theme-provider.component';
import { ProfileSkeletonComponent } from './skeletons/dashboard/profile/profile-skeleton.component';

@NgModule({
  imports: [CommonModule],
  declarations: [
    ThemeProviderComponent,
    EmailSentComponent,
    ProfileSkeletonComponent,
  ],
  exports: [
    ThemeProviderComponent,
    EmailSentComponent,
    ProfileSkeletonComponent,
  ],
})
export class ComponentsModule {}
