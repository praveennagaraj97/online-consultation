import { Component, Input } from '@angular/core';
import { DoctorsEntity } from 'src/app/types/doctor.response.types';

@Component({
  selector: 'app-doctor-list-data',
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
          <th class="text-sm font-medium p-2 whitespace-nowrap text-center">
            Status
          </th>
          <th class="text-sm font-medium p-2 whitespace-nowrap text-center">
            Actions
          </th>
        </thead>
        <tbody *ngIf="doctors.length; else notfound">
          <tr
            *ngFor="let doctor of doctors"
            class="smooth-animate hover:bg-gray-400/30 border-b dark:border-gray-50/20"
          >
            <td class="text-xs p-2 whitespace-nowrap">
              <div class="flex items-center space-x-1">
                <img
                  [src]="
                    doctor.profile_pic?.blur_data_url ||
                    '/assets/img-placeholder.png'
                  "
                  width="32"
                  height="32"
                  alt="img"
                  class="rounded-full w-8 h-8 min-w-[32px] min-h-[32px] overflow-hidden"
                />
                <div class="flex flex-col">
                  <span>
                    {{ doctor.name }}
                  </span>
                  <small>
                    {{ doctor.professional_title }}
                  </small>
                </div>
              </div>
            </td>
            <td class="text-xs p-2 whitespace-nowrap">
              {{ doctor.email }}
            </td>
            <td class="text-xs p-2 whitespace-nowrap">
              {{ doctor.phone.code + ' ' + doctor.phone.number }}
            </td>
            <td class="text-xs p-2 whitespace-nowrap">
              {{ doctor.education }}
            </td>
            <td class="text-xs p-2 whitespace-nowrap">
              {{ doctor.experience }}
              {{ doctor.experience > 1 ? 'years' : 'year' }}
            </td>
            <td class="text-xs p-2 whitespace-nowrap">
              {{ doctor?.hospital?.name }}
            </td>
            <td class="text-xs p-2 whitespace-nowrap">
              <select
                aria-readonly="true"
                name="languages"
                class="common-input w-full input-focus input-colors rounded-md p-1"
              >
                <option
                  aria-readonly="true"
                  *ngFor="let lang of doctor.spoken_languages"
                >
                  {{ lang.name }}
                </option>
              </select>
            </td>
            <td class="text-xs p-2 whitespace-nowrap">
              <div class="flex justify-center">
                <app-doctor-status
                  [isActive]="doctor.is_active"
                ></app-doctor-status>
              </div>
            </td>
            <td class="text-xs p-2 whitespace-nowrap">
              <app-three-dots-icon
                className="h-5 w-5  fill-gray-500 stroke-gray-500 dark:fill-gray-100 dark:stroke-gray-100 mx-auto"
              ></app-three-dots-icon>
            </td>
          </tr>
        </tbody>

        <ng-template #notfound>
          <tbody>
            <tr>
              <td colspan="9">
                <div class="flex items-center justify-center">
                  <app-no-doctors></app-no-doctors>
                </div>
              </td>
            </tr>
          </tbody>
        </ng-template>
      </table>
    </div>
  `,
})
export class DoctorListTableDataComponent {
  @Input() doctors: DoctorsEntity[] = [];
}
