import { Component, Input } from '@angular/core';
import { DoctorsEntity } from 'src/app/types/doctor.response.types';

@Component({
  selector: 'app-doctor-list-data-component',
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
          <th class="text-sm font-medium p-2 whitespace-nowrap text-left">
            Actions
          </th>
        </thead>
        <tbody>
          <tr
            *ngFor="let doctor of doctors"
            class="smooth-animate hover:bg-gray-400/30"
          >
            <td class="text-xs p-2 whitespace-nowrap">
              <div class="flex items-center space-x-1">
                <img
                  [src]="
                    doctor.profile_pic?.image_src ||
                    '/assets/img-placeholder.png'
                  "
                  width="48"
                  height="48"
                  alt="img"
                  class="rounded-full w-8 h-8 overflow-hidden"
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
                class="common-input input-focus input-colors rounded-md p-1"
              >
                <option value="" *ngFor="let lang of doctor.spoken_languages">
                  {{ lang.name }}
                </option>
              </select>
            </td>
            <td class="text-xs p-2 whitespace-nowrap">
              {{ doctor.is_active ? 'active' : 'inactive' }}
            </td>
            <td class="text-xs p-2 whitespace-nowrap">...</td>
          </tr>
        </tbody>
      </table>
    </div>
  `,
})
export class DoctorListTableDataComponent {
  @Input() doctors: DoctorsEntity[] = [];
}
