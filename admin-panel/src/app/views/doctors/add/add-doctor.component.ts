import { transition, trigger, useAnimation } from '@angular/animations';
import { Component } from '@angular/core';
import {
  FormBuilder,
  FormControl,
  FormGroup,
  Validators,
} from '@angular/forms';
import { Subscription } from 'rxjs';
import { fadeInTransformAnimation } from 'src/app/animations';
import { doctorFormErrors } from 'src/app/errors/doctor-form.errors';
import { BreadcrumbPath, SelectOption } from 'src/app/types/app.types';
import { clearSubscriptions } from 'src/app/utils/helpers';
import { AddDoctorService } from './add-doctor.service';

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
  // Subs
  private subs$: Subscription[] = [];

  constructor(
    private fb: FormBuilder,
    private addNewDocService: AddDoctorService
  ) {}

  // State
  breadcrumbPaths: BreadcrumbPath[] = [
    { path: '/doctors', title: 'Doctors' },
    { path: '/doctors/add', title: 'Add New Doctor' },
  ];
  shouldShowError = false;
  errors = doctorFormErrors;
  hospitalOptions: SelectOption[] = [];
  hospitalsLoading = false;
  nextHospitalsPaginateId: string | null = null;
  hospitalSearchTerm = '';

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
    hospital_id: new FormControl(''),
    // Ref
    hospital_title: new FormControl(''),
  });

  ngOnInit() {
    this.getHospitalsOptions();
  }

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

  private getHospitalsOptions(shouldReset = false) {
    this.hospitalsLoading = true;
    this.subs$.push(
      this.addNewDocService
        .getHospitals(
          this.nextHospitalsPaginateId,
          this.hospitalSearchTerm,
          shouldReset
        )
        .subscribe({
          next: (res) => {
            this.hospitalsLoading = false;
            this.hospitalOptions = res.hospitals;
            this.nextHospitalsPaginateId = res.nextId;
          },
          error: (err) => {
            this.hospitalsLoading = false;
            alert('Failed to load hospitals');
          },
        })
    );
  }

  loadMoreHospitals() {
    this.getHospitalsOptions();
  }

  onHospitalSearch(term: string) {
    this.nextHospitalsPaginateId = null;
    this.hospitalsLoading = true;
    this.hospitalOptions = [];
    this.hospitalSearchTerm = term;
    this.getHospitalsOptions(true);
  }

  onHospitalSelect(opt: SelectOption) {
    this.doctorForm.controls?.['hospital_title'].setValue(opt.title);
    this.doctorForm.controls?.['hospital_id'].setValue(opt.value);
  }

  ngOnDestroy() {
    clearSubscriptions(this.subs$);
  }
}
