import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-pagination-count-display-component',
  template: `
    <small class="text-xs block" *ngIf="!loading; else skeleton">
      Showing {{ currentCount }} of {{ totalCount }}
    </small>
    <ng-template #skeleton>
      <div class="flex space-x-2">
        <div class="skeleton w-16 h-2 rounded-md"></div>
        <div class="skeleton w-4 h-2 rounded-md"></div>
        <div class="skeleton w-4 h-2 rounded-md"></div>
        <div class="skeleton w-6 h-2 rounded-md"></div>
      </div>
    </ng-template>
  `,
})
export class PaginationCountDisplayComponent {
  @Input() totalCount = 0;
  @Input() currentCount = 0;
  @Input() loading = true;
}
