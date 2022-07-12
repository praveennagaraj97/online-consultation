import { Component } from '@angular/core';

@Component({
  selector: 'app-pagination-controls-skeleton-component',
  template: `
    <div class="flex items-center space-x-1">
      <app-previous-icon
        className="w-6 h-6 animate-pulse dark:fill-gray-50 dark:stroke-gray-50"
      ></app-previous-icon>

      <app-next-icon
        className="w-6 h-6 animate-pulse dark:fill-gray-50 dark:stroke-gray-50"
      ></app-next-icon>
    </div>
  `,
})
export class PaginationControlsSkeletonComponent {}
