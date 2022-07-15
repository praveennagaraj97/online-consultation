import { Component } from '@angular/core';

@Component({
  selector: 'app-pagination-count-skeleton',
  template: `
    <div class="flex space-x-2">
      <div class="skeleton w-16 h-2 rounded-md"></div>
      <div class="skeleton w-4 h-2 rounded-md"></div>
      <div class="skeleton w-4 h-2 rounded-md"></div>
      <div class="skeleton w-6 h-2 rounded-md"></div>
    </div>
  `,
})
export class PaginationSkeletonComponent {}
