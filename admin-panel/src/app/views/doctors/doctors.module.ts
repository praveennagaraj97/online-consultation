import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';
import { ComponentsModule } from 'src/app/components/component.module';
import { AddNewDoctorViewComponent } from './add/add-doctor.component';
import { AddDoctorService } from './add/add-doctor.service';
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
  {
    path: 'add',
    component: AddNewDoctorViewComponent,
    title: 'Online Consultation | Doctor | New',
  },
];

@NgModule({
  declarations: [DoctorsListViewComponent, AddNewDoctorViewComponent],
  imports: [
    CommonModule,
    RouterModule.forChild(routes),
    ComponentsModule,
    FormsModule,
    ReactiveFormsModule,
  ],
  providers: [DoctorsListViewService, AddDoctorService],
})
export class DoctorsViewModule {}
