import { transition, trigger, useAnimation } from '@angular/animations';
import { Component } from '@angular/core';
import { Subscription } from 'rxjs';
import { fadeInTransformAnimation } from 'src/app/animations';
import { SelectOption } from 'src/app/types/app.types';
import { DoctorsEntity } from 'src/app/types/doctor.response.types';
import { clearSubscriptions } from 'src/app/utils/helpers';
import { DoctorsListViewService } from './list-view.service';

@Component({
  selector: 'app-doctors-list-view',
  templateUrl: 'list-view.component.html',
  animations: [
    trigger('fadeInOut', [
      transition('void <=> *', [useAnimation(fadeInTransformAnimation())]),
    ]),
  ],
})
export class DoctorsListViewComponent {
  // Subs
  private subs$: Subscription[] = [];

  // State
  breadcrumbPaths = [{ path: '/doctors', title: 'Doctors' }];
  isLoading = false;
  doctors: DoctorsEntity[] = [];
  perPage = 10;
  specialities: SelectOption[] = [];
  consultationTypes: SelectOption[] = [];
  activeStatus: SelectOption[] = [
    { title: 'All', value: '' },
    { title: 'Active', value: 'true' },
    { title: 'In Active', value: 'false' },
  ];
  selectedSpeciality = '';
  selectedConsultationType = '';
  selectedActiveType = '';

  constructor(private doctorsListService: DoctorsListViewService) {}

  ngOnInit() {
    this.getDoctorsList();
    this.getSpecialities();
    this.getConsultationTypes();
  }

  private getDoctorsList() {
    this.isLoading = true;
    this.subs$.push(
      this.doctorsListService
        .getDoctorsList(
          this.perPage,
          this.selectedSpeciality,
          this.selectedConsultationType,
          this.selectedActiveType
        )
        .subscribe({
          next: (res) => {
            this.isLoading = false;
            this.doctors = res.result || [];
          },
          error: () => {
            this.isLoading = false;
            alert('Failed to load doctors');
          },
        })
    );
  }

  private getSpecialities() {
    this.subs$.push(
      this.doctorsListService.getSpecialities().subscribe({
        next: (res) => {
          this.specialities = res || [];
        },
        error: () => {
          alert('Failed to load specialities');
        },
      })
    );
  }

  private getConsultationTypes() {
    this.subs$.push(
      this.doctorsListService.getConsultationTypes().subscribe({
        next: (res) => {
          this.consultationTypes = res || [];
        },
        error: () => {
          alert('Failed to load specialities');
        },
      })
    );
  }

  onSpecialityChanegFilter(option: SelectOption) {
    this.selectedSpeciality = option.value;
    this.getDoctorsList();
  }

  onConsultationChangeFilter(option: SelectOption) {
    this.selectedConsultationType = option.value;
    this.getDoctorsList();
  }

  onActiveStatusChangeFilter(option: SelectOption) {
    this.selectedActiveType = option.value;
    this.getDoctorsList();
  }

  updatePerPageLimit(value: number) {
    this.perPage = value;

    this.getDoctorsList();
  }

  ngOnDestropy() {
    clearSubscriptions(this.subs$);
  }
}
