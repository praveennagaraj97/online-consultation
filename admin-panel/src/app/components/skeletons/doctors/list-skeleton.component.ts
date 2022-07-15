import { Component } from '@angular/core';

@Component({
  selector: 'app-doctors-list-skeleton',
  template: `
    <div class="overflow-x-auto inner-scrollbar">
      <table class="w-full min-w-[1000px]">
        <thead class="w-full bg-gray-400/50">
          <th class="text-sm font-medium p-2 whitespace-nowrap text-left">
            Name
          </th>
          <th class="text-sm font-medium p-2 whitespace-nowrap text-left">
            Email
          </th>
          <th class="text-sm font-medium p-2 whitespace-nowrap text-left">
            Phone
          </th>
          <th class="text-sm font-medium p-2 whitespace-nowrap text-left">
            Education
          </th>
          <th class="text-sm font-medium p-2 whitespace-nowrap text-left">
            Experience
          </th>
          <th class="text-sm font-medium p-2 whitespace-nowrap text-left">
            Hospital
          </th>
          <th class="text-sm font-medium p-2 whitespace-nowrap text-left">
            languages
          </th>
          <th class="text-sm font-medium p-2 whitespace-nowrap text-left">
            Status
          </th>
          <th class="text-sm font-medium p-2 whitespace-nowrap text-center">
            Actions
          </th>
        </thead>
        <tbody>
          <tr
            *ngFor="let _ of repeatArray"
            class="border-b dark:border-gray-50/20"
          >
            <td class="p-2">
              <div class="flex items-center space-x-1">
                <div class="w-8 h-8 rounded-full skeleton"></div>
                <div class="flex flex-col">
                  <div class="skeleton w-24 h-3 rounded-md mb-1"></div>
                  <div class="skeleton w-12 h-2 rounded-md"></div>
                </div>
              </div>
            </td>
            <td class="p-2">
              <div class="skeleton w-10/12 h-3 rounded-md"></div>
            </td>
            <td class="p-2">
              <div class="skeleton w-10/12 h-3 rounded-md"></div>
            </td>
            <td class="p-2">
              <div class="skeleton w-10/12 h-3 rounded-md"></div>
            </td>
            <td class="p-2">
              <div class="skeleton w-10/12 h-3 rounded-md"></div>
            </td>
            <td class="p-2">
              <div class="skeleton w-10/12 h-3 rounded-md"></div>
            </td>
            <td class="p-2">
              <div class="skeleton w-10/12 h-3 rounded-md"></div>
            </td>
            <td class="p-2">
              <div class="skeleton w-10/12 h-3 rounded-md"></div>
            </td>
            <td class="p-2 h-full">
              <app-three-dots-icon
                className="h-5 w-5  fill-gray-300 stroke-gray-300 dark:fill-gray-500 dark:stroke-gray-500 mx-auto animate-pulse"
              ></app-three-dots-icon>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  `,
})
export class DoctorsListSkeletonComponent {
  repeatArray: ReadonlyArray<[]> = new Array(10).fill('');
}
