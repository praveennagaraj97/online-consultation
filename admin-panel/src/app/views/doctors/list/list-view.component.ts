import { transition, trigger, useAnimation } from '@angular/animations';
import { Component } from '@angular/core';
import { fadeInTransformAnimation } from 'src/app/animations';
import { DoctorsEntity } from 'src/app/types/doctor.response.types';
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
  // State
  breadcrumbPaths = [{ path: '/doctors', title: 'Doctors' }];
  doctors: DoctorsEntity[] = [];

  constructor(private doctorsListService: DoctorsListViewService) {}

  ngOnInit() {
    this.doctorsListService.getDoctorsList().subscribe({
      next: (res) => {
        console.log(res);
        this.doctors = res.result || [];
      },
      error: (err) => {
        console.log(err);
      },
    });
  }
}
