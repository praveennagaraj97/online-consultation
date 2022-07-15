import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-pagination-count-display',
  template: `
    <small class="text-xs block">
      Showing {{ startCount }} - {{ endCount }} of {{ totalCount }}
    </small>
  `,
})
export class PaginationCountDisplayComponent {
  @Input() totalCount = 0;
  @Input() startCount = 0;
  @Input() endCount = 0;
}
