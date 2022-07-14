import { transition, trigger, useAnimation } from '@angular/animations';
import { Component } from '@angular/core';
import {
  FormArray,
  FormBuilder,
  FormControl,
  FormGroup,
  Validators,
} from '@angular/forms';
import { Subscription } from 'rxjs';
import { fadeInTransformAnimation } from 'src/app/animations';
import { doctorFormErrors } from 'src/app/errors/doctor-form.errors';
import { BreadcrumbPath, SelectOption } from 'src/app/types/app.types';
import type { DoctorFormDTO } from 'src/app/types/dto.types';
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
  activeStatusOptions: SelectOption[] = [
    { title: 'Active', value: 'true' },
    { title: 'In Active', value: 'false' },
  ];
  shouldShowError = false;
  errors = doctorFormErrors;
  hospitalOptions: SelectOption[] = [];
  hospitalsLoading = false;
  nextHospitalsPaginateId: string | null = null;
  hospitalSearchTerm = '';
  consultationTypeOptions: SelectOption[] = [];
  specialityOptions: SelectOption[] = [];
  specialityHasNext: string | null = null;
  specialityLoading = false;
  specialitySearchTerm = '';
  languagesOptions: SelectOption[] = [];
  languagesHasNext: string | null = null;
  languagesLoading = false;
  languagesSearchTerm = '';

  // Input Form Group State
  profilePic: File | null = null;
  doctorForm: FormGroup<DoctorFormDTO> = this.fb.group({
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
    consultation_type_id: new FormControl('', {
      validators: [Validators.required],
    }),
    speciality_id: new FormControl(''),
    spoken_language_id: new FormArray<FormControl<string>>([], {
      validators: [Validators.required, Validators.min(1)],
    }),
    is_active: new FormControl('false'),
  });

  ngOnInit() {
    this.getHospitalsOptions();
    this.getConsultationTypes();
    this.getSpecialities();
    this.getLanguages();
  }

  getFormValue(name: string): string {
    if (name == 'name' && !this.doctorForm.get(name)?.dirty) {
      return '';
    }

    return this.doctorForm.get(name)?.value || '';
  }

  onFormSubmit(form: FormGroup<DoctorFormDTO>) {
    if (form.invalid) {
      this.shouldShowError = true;
    } else {
      this.addNewDocService.handleAddDoctor(form, this.profilePic).subscribe({
        next: (val) => {
          console.log(val);
        },
        error: (err) => {
          console.log(err);
        },
      });
    }
  }

  // Hospital input
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
          error: () => {
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
    this.doctorForm.controls?.['hospital_id'].setValue(opt.value);
  }

  // Consultation Type Input
  private getConsultationTypes() {
    this.subs$.push(
      this.addNewDocService.getConsultaionTypeOptions().subscribe({
        next: (options) => {
          this.consultationTypeOptions = options || [];
        },
        error: (err) => {
          alert(err);
        },
      })
    );
  }

  onConsultationTypeSelect(id: string) {
    this.doctorForm.controls?.['consultation_type_id']?.setValue(id);
  }

  // Specialities
  private getSpecialities(shouldReset = false) {
    this.specialityLoading = true;
    this.subs$.push(
      this.addNewDocService
        .getSpecialitiesOptions(
          this.nextHospitalsPaginateId,
          this.specialitySearchTerm,
          shouldReset
        )
        .subscribe({
          next: (options) => {
            this.specialityOptions = options.specialities || [];
            this.specialityHasNext = options.nextId;
            this.specialityLoading = false;
          },
          error: (err) => {
            this.specialityLoading = false;
            alert(err);
          },
        })
    );
  }

  loadMoreSpecialities() {
    this.getSpecialities();
  }

  onSpecialitySearch(term: string) {
    this.nextHospitalsPaginateId = null;
    this.specialityLoading = true;
    this.specialityOptions = [];
    this.specialitySearchTerm = term;
    this.getSpecialities(true);
  }

  onSpecialitySelect(opt: SelectOption) {
    this.doctorForm.controls?.['speciality_id'].setValue(opt.value);
  }

  // Languages
  private getLanguages(shouldReset = false) {
    this.languagesLoading = true;
    this.subs$.push(
      this.addNewDocService
        .getLangaugesOptions(
          this.nextHospitalsPaginateId,
          this.languagesSearchTerm,
          shouldReset
        )
        .subscribe({
          next: (options) => {
            this.languagesOptions = options.langauges || [];
            this.languagesHasNext = options.nextId;
            this.languagesLoading = false;
          },
          error: (err) => {
            this.languagesLoading = false;
            alert(err);
          },
        })
    );
  }

  loadMoreLanguges() {
    this.getLanguages();
  }

  onLanguageSearch(term: string) {
    this.nextHospitalsPaginateId = null;
    this.languagesLoading = true;
    this.languagesOptions = [];
    this.languagesSearchTerm = term;
    this.getLanguages(true);
  }

  onLanguageSelect(opt: SelectOption) {
    const input = this.doctorForm.controls?.['spoken_language_id'] as FormArray;

    if (input.value?.includes(opt.value)) {
      input.removeAt(input.value?.findIndex((val: string) => val == opt.value));
      return;
    }
    if (input) {
      input.push(new FormControl(opt.value));
    }
  }

  get checkIfConsultationTypeIsSchedule() {
    return (
      this.consultationTypeOptions.find(
        (type) =>
          type.value ==
          this.doctorForm.controls?.['consultation_type_id']?.value
      )?.title === 'Schedule'
    );
  }

  handleProfilePicChange(files: File[]) {
    this.profilePic = files[0];
  }

  ngOnDestroy() {
    clearSubscriptions(this.subs$);
  }
}
