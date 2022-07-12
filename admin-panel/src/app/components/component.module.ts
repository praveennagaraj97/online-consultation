import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { IconsModule } from '../icons/icons.module';
import { BreadcrumbComponent } from './breadcrumb/breadcrumb.component';
import { DoctorListTableDataComponent } from './doctors/table-data/table-data.component';
import { EmailSentComponent } from './email-sent/email-sent.component';
import { PerPageOptionsComponent } from './pagination/per-page-options.component';
import { ThemeProviderComponent } from './providers/theme-provider/theme-provider.component';
import { ProfileSkeletonComponent } from './skeletons/dashboard/profile/profile-skeleton.component';

@NgModule({
  imports: [CommonModule, IconsModule, RouterModule, FormsModule],
  declarations: [
    ThemeProviderComponent,
    EmailSentComponent,
    ProfileSkeletonComponent,
    BreadcrumbComponent,
    PerPageOptionsComponent,
    DoctorListTableDataComponent,
  ],
  exports: [
    ThemeProviderComponent,
    EmailSentComponent,
    ProfileSkeletonComponent,
    BreadcrumbComponent,
    PerPageOptionsComponent,
    DoctorListTableDataComponent,
  ],
})
export class ComponentsModule {}
