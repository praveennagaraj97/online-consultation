import {
  animate,
  state,
  style,
  transition,
  trigger,
} from '@angular/animations';
import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-toggle-input',
  template: `
    <div class="flex items-center space-x-2 ">
      <button
        class="w-8 h-5 flex items-center rounded-full"
        [ngClass]="isActive ? 'bg-green-500' : 'bg-red-500'"
        (click)="onToggle.emit(!isActive)"
        [disabled]="isLoading"
      >
        <button
          class="w-4 h-4 rounded-full bg-gray-50 mx-1 flex items-center justify-center"
          [@toggle]="isActive ? 'active' : 'inactive'"
        >
          <svg
            className="w-4 h-4 fill-gray-700 stroke-gray-700"
            viewBox="0 0 1024 1024"
            xmlns="http://www.w3.org/2000/svg"
            *ngIf="!isActive && !isLoading"
          >
            <path
              d="M696 480H328c-4.4 0-8 3.6-8 8v48c0 4.4 3.6 8 8 8h368c4.4 0 8-3.6 8-8v-48c0-4.4-3.6-8-8-8z"
            ></path>
            <path
              d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z"
            ></path>
          </svg>
          <svg
            className="w-4 h-4 fill-gray-700 stroke-gray-700"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
            *ngIf="isActive && !isLoading"
          >
            <path
              fill="none"
              stroke="#000"
              stroke-width="2"
              d="M12,22 C17.5228475,22 22,17.5228475 22,12 C22,6.4771525 17.5228475,2 12,2 C6.4771525,2 2,6.4771525 2,12 C2,17.5228475 6.4771525,22 12,22 Z M7,12 L11,15 L16,8"
            ></path>
          </svg>
          <app-spinner-icon
            *ngIf="isLoading"
            className="w-3.5 h-3.5 fill-gray-700 stroke-gray-700 animate-spin"
          ></app-spinner-icon>
        </button>
      </button>
    </div>
  `,
  animations: [
    trigger('toggle', [
      state('inactive', style({ transform: 'translateX(-1.5px)' })),
      state('active', style({ transform: 'translateX(10px)' })),
      transition('active <=> inactive', [animate('0.4s')]),
    ]),
  ],
})
export class ToggleInputComponent {
  @Input() isActive = false;
  @Input() isLoading = false;
  @Output() onToggle = new EventEmitter<boolean>();
}
