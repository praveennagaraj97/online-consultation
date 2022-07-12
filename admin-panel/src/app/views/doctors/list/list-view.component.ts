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
  doctors: DoctorsEntity[] = [];
  perPage = 10;
  limitOptions: { value: number; title: string }[] = [
    { title: '10', value: 10 },
    { title: '20', value: 20 },
    { title: '50', value: 50 },
  ];
  specialities: SelectOption[] = [];
  selectedSpeciality = '';

  constructor(private doctorsListService: DoctorsListViewService) {}

  ngOnInit() {
    this.getDoctorsList();
    this.getSpecialities();
  }

  private getSpecialities() {
    this.subs$.push(
      this.doctorsListService.getSpecialities().subscribe({
        next: (res) => {
          this.specialities = res || [];
        },
        error: (err) => {
          console.log(err);
        },
      })
    );
  }

  // Error Handle Pending
  private getDoctorsList() {
    this.subs$.push(
      this.doctorsListService
        .getDoctorsList(this.perPage, this.selectedSpeciality)
        .subscribe({
          next: (res) => {
            this.doctors = res.result || [];
          },
          error: (err) => {
            console.log(err);
          },
        })
    );
  }

  onSpecialityFilter(option: SelectOption) {
    this.selectedSpeciality = option.value;
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
