import { Component } from '@angular/core';

@Component({
  selector: 'app-per-page-component',
  template: `
    <div class="flex items-center space-x-1.5">
      <small>Show</small>
      <select
        name="per_page"
        class="common-input input-focus input-colors rounded-md p-1 text-sm"
      >
        <option value="">10</option>
        <option value="">20</option>
        <option value="">50</option>
      </select>
      <small>Entries</small>
    </div>
  `,
})
export class PerPageOptionsComponent {}
