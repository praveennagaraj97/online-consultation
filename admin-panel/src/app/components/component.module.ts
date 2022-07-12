import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { DirectiveModule } from '../directives/directive.module';
import { IconsModule } from '../icons/icons.module';
import { NoDoctorsFoundComponent } from './404s/no-doctors.component';
import { BreadcrumbComponent } from './breadcrumb/breadcrumb.component';
import { DoctorListTableDataComponent } from './doctors/table-data/table-data.component';
import { EmailSentComponent } from './email-sent/email-sent.component';
import { PaginationControlsComponent } from './pagination/controls.component';
import { PaginationCountDisplayComponent } from './pagination/display-count.component';
import { PerPageOptionsComponent } from './pagination/per-page-options.component';
import { ThemeProviderComponent } from './providers/theme-provider/theme-provider.component';
import { SelectDropdownComponent } from './select-dropdown/select-dropdown.component';
import { ProfileSkeletonComponent } from './skeletons/dashboard/profile/profile-skeleton.component';
import { DoctorsListSkeletonComponent } from './skeletons/doctors/list-skeleton.component';

@NgModule({
  imports: [
    CommonModule,
    IconsModule,
    RouterModule,
    FormsModule,
    DirectiveModule,
  ],
  declarations: [
    ThemeProviderComponent,
    EmailSentComponent,
    ProfileSkeletonComponent,
    BreadcrumbComponent,
    PerPageOptionsComponent,
    DoctorListTableDataComponent,
    SelectDropdownComponent,
    NoDoctorsFoundComponent,
    DoctorsListSkeletonComponent,
    PaginationControlsComponent,
    PaginationCountDisplayComponent,
  ],
  exports: [
    ThemeProviderComponent,
    EmailSentComponent,
    ProfileSkeletonComponent,
    BreadcrumbComponent,
    PerPageOptionsComponent,
    DoctorListTableDataComponent,
    SelectDropdownComponent,
    NoDoctorsFoundComponent,
    DoctorsListSkeletonComponent,
    PaginationControlsComponent,
    PaginationCountDisplayComponent,
  ],
})
export class ComponentsModule {}
