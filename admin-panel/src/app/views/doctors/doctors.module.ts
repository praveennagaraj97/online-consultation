import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ComponentsModule } from 'src/app/components/component.module';
import { EditDoctorViewComponent } from './edit/edit-doctor.component';
import { DoctorsListViewComponent } from './list/list-view.component';
import { DoctorsListViewService } from './list/list-view.service';

const routes: Routes = [
  {
    path: '',
    component: DoctorsListViewComponent,
    title: 'Online Consultation | Doctors',
  },
  {
    path: 'edit/:id',
    component: EditDoctorViewComponent,
    title: 'Online Consultation | Doctor | Edit',
  },
];

@NgModule({
  declarations: [DoctorsListViewComponent],
  imports: [CommonModule, RouterModule.forChild(routes), ComponentsModule],
  providers: [DoctorsListViewService],
})
export class DoctorsViewModule {}
