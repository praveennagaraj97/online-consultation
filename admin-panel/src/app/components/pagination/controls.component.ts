import { Location } from '@angular/common';
import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-pagination-controls',
  template: `
    <div class="flex items-center space-x-1">
      <button
        class="smooth-animate hover:scale-110"
        [disabled]="pageNum == 1"
        (click)="goBack()"
      >
        <app-previous-icon
          className="w-6 h-6 dark:fill-gray-50 dark:stroke-gray-50"
        ></app-previous-icon>
      </button>
      <button
        class="smooth-animate hover:scale-110"
        [disabled]="!paginateId"
        (click)="goNext()"
      >
        <a
          [routerLink]="nextPageLink"
          [queryParams]="{ paginate_cur: paginateId, page_num: pageNum + 1 }"
        >
          <app-next-icon
            className="w-6 h-6 dark:fill-gray-50 dark:stroke-gray-50"
          ></app-next-icon>
        </a>
      </button>
    </div>
  `,
})
export class PaginationControlsComponent {
  @Input() nextPageLink = '';
  @Input() pageNum = 1;
  @Input() paginateId = '';

  constructor(private history: Location) {}

  goNext() {}

  goBack() {
    this.history.back();
  }
}
