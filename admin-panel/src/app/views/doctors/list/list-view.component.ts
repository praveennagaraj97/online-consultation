import { transition, trigger, useAnimation } from '@angular/animations';
import { Component } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
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
  totalCount = 0;
  pageNum = 1;
  paginateId = '';
  searchTerm = '';

  constructor(
    private doctorsListService: DoctorsListViewService,
    private activeRoute: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit() {
    this.getDoctorsList();
    this.getSpecialities();
    this.getConsultationTypes();
    this.listenToBrowserHistory();
  }

  private getDoctorsList() {
    const url = this.router.url;
    const paginateId =
      this.router.parseUrl(url).queryParams?.['paginate_cur'] || '';
    this.pageNum =
      parseInt(this.router.parseUrl(url).queryParams?.['page_num']) || 1;

    this.isLoading = true;
    this.subs$.push(
      this.doctorsListService
        .getDoctorsList(
          this.perPage,
          this.selectedSpeciality,
          this.selectedConsultationType,
          this.selectedActiveType,
          paginateId,
          this.searchTerm
        )
        .subscribe({
          next: (res) => {
            this.isLoading = false;
            this.doctors = res.result || [];
            this.totalCount = res.count;
            this.paginateId = res.paginate_id || '';
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

  onSpecialityChangeFilter(option: SelectOption) {
    this.router.navigate([], { queryParams: {} });

    this.selectedSpeciality = option.value;
    this.getDoctorsList();
  }

  onConsultationChangeFilter(option: SelectOption) {
    this.router.navigate([], { queryParams: {} });
    this.selectedConsultationType = option.value;
    this.getDoctorsList();
  }

  onActiveStatusChangeFilter(option: SelectOption) {
    this.router.navigate([], { queryParams: {} });
    this.selectedActiveType = option.value;
    this.getDoctorsList();
  }

  updatePerPageLimit(value: number) {
    this.router.navigate([], { queryParams: {} });
    this.perPage = value;
    this.getDoctorsList();
  }

  onSearch(value: string) {
    this.router.navigate([], { queryParams: {} });
    this.searchTerm = value;
    this.getDoctorsList();
  }

  private listenToBrowserHistory() {
    this.activeRoute.queryParams.subscribe({
      next: () => {
        this.getDoctorsList();
      },
    });
  }

  ngOnDestropy() {
    clearSubscriptions(this.subs$);
  }
}
