import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { AdminIconComponent } from './admin/admin-icon.component';
import { AppointmentIconComponent } from './appointment/appointment-icon.component';
import { ConsultatationIconComponent } from './consultation/consultation-icon.component';
import { DoctorIconComponent } from './doctor/doctor-icon.component';
import { EditorIconComponent } from './editor/editor-icon.component';
import { HomeIconComponent } from './home/home-icon.component';
import { HospitalIconComponent } from './hospital/hospital-icon.component';
import { LanguageIconComponent } from './language/language-icon.component';
import { LockIconComponent } from './lock/lock-icon.component';
import { MailIconComponent } from './mail/mail-icon.component';
import { SpecialityIconComponent } from './speciality/speciality-icon.component';
import { SpinnerIconComponent } from './spinner/spinner-icon.component';
import { UserIconComponent } from './user/user-icon.component';

@NgModule({
  imports: [CommonModule],
  declarations: [
    MailIconComponent,
    LockIconComponent,
    SpinnerIconComponent,
    HomeIconComponent,
    DoctorIconComponent,
    UserIconComponent,
    AppointmentIconComponent,
    AdminIconComponent,
    EditorIconComponent,
    HospitalIconComponent,
    ConsultatationIconComponent,
    LanguageIconComponent,
    SpecialityIconComponent,
  ],
  exports: [
    MailIconComponent,
    LockIconComponent,
    SpinnerIconComponent,
    HomeIconComponent,
    DoctorIconComponent,
    UserIconComponent,
    AppointmentIconComponent,
    AdminIconComponent,
    EditorIconComponent,
    HospitalIconComponent,
    ConsultatationIconComponent,
    LanguageIconComponent,
    SpecialityIconComponent,
  ],
})
export class IconsModule {}
