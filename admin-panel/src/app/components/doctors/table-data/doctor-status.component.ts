import { Component, Input } from '@angular/core';
import { ConfirmPortalEventTypes } from 'src/app/types/app.types';

@Component({
  selector: 'app-doctor-status',
  template: `<app-toggle-input
      [title]="
        isActive ? 'click to mark as inactive' : 'click to mark as active'
      "
      [isActive]="isActive"
      (onToggle)="showConfirmModal = true"
      [isLoading]="isLoading"
      stopPropagation
    ></app-toggle-input>

    <app-confirm-dialog-portal
      (onConfirm)="onConfirm($event)"
      [showModal]="showConfirmModal"
    ></app-confirm-dialog-portal> `,
})
export class DoctorStatusToggleComponent {
  @Input() isActive = false;
  isLoading = false;

  showConfirmModal = false;

  updateDoctorStatus(status: boolean) {
    this.isLoading = true;
  }

  onConfirm(reason: ConfirmPortalEventTypes) {
    if (reason == 'cancel') {
      this.showConfirmModal = false;
    }
  }
}
