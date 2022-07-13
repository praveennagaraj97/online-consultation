import { transition, trigger, useAnimation } from '@angular/animations';
import { Component } from '@angular/core';
import {
  FormBuilder,
  FormControl,
  FormGroup,
  Validators,
} from '@angular/forms';
import { fadeInTransformAnimation } from 'src/app/animations';
import { doctorFormErrors } from 'src/app/errors/doctor-form.errors';
import { BreadcrumbPath } from 'src/app/types/app.types';

@Component({
  selector: 'app-add-new-doctor-view',
  templateUrl: 'add-doctor.component.html',
  animations: [
    trigger('fadeInOut', [
      transition('void <=> *', [useAnimation(fadeInTransformAnimation())]),
    ]),
    trigger('fadeIn', [
      transition('void => *', [useAnimation(fadeInTransformAnimation())]),
    ]),
  ],
})
export class AddNewDoctorViewComponent {
  constructor(private fb: FormBuilder) {}

  // State
  breadcrumbPaths: BreadcrumbPath[] = [
    { path: '/doctors', title: 'Doctors' },
    { path: '/doctors/add', title: 'Add New Doctor' },
  ];
  shouldShowError = false;
  errors = doctorFormErrors;

  // Input Form Group State
  profilePic: File | null = null;
  doctorForm: FormGroup = this.fb.group({
    name: new FormControl('Dr. ', {
      validators: [Validators.required],
    }),
    email: new FormControl('', {
      validators: [Validators.required, Validators.email],
      asyncValidators: [],
    }),
    professional_title: new FormControl('', {
      validators: [Validators.required],
    }),
    phone_code: new FormControl('+91'),
    phone_number: new FormControl('', {
      validators: [Validators.required, Validators.pattern('^[0-9]+$')],
    }),
    education: new FormControl('', {
      validators: [Validators.required],
    }),
    experience: new FormControl('', {
      validators: [Validators.required, Validators.pattern('^[0-9]+$')],
    }),
  });

  getFormValue(name: string): string {
    if (name == 'name' && !this.doctorForm.get(name)?.dirty) {
      return '';
    }

    return this.doctorForm.get(name)?.value || '';
  }

  onFormSubmit(form: FormGroup) {
    if (form.invalid) {
      this.shouldShowError = true;
    } else {
      console.log(form.value);
    }
  }
}
