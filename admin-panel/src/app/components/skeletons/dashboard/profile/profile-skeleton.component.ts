import { Component } from '@angular/core';

@Component({
  selector: 'app-profile-skeleton',
  template: `
    <div class="flex items-center">
      <div class="mr-2">
        <div class="flex space-x-1">
          <div class="shadow-sm w-20 h-3 rounded-lg skeleton mb-1"></div>
          <div class="shadow-sm w-10 h-3 rounded-lg skeleton mb-1"></div>
        </div>
        <div class="shadow-sm w-10 h-2 rounded-lg skeleton ml-auto"></div>
      </div>
      <div class="skeleton rounded-full w-10 h-10 shadow-sm"></div>
    </div>
  `,
})
export class ProfileSkeletonComponent {}
