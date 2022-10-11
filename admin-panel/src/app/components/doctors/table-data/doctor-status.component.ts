import { Component, Input } from '@angular/core';
import { ConfirmPortalEventTypes } from 'src/app/types/app.types';
import { ConfirmDialogPortalService } from '../../portal/dialogs/confirm/confirm-dialog-portal.service';

@Component({
  selector: 'app-doctor-status',
  template: `<app-toggle-input
      [title]="
        isActive ? 'click to mark as inactive' : 'click to mark as active'
      "
      [isActive]="isActive"
      (onToggle)="showConfirmModal = true"
      stopPropagation
    ></app-toggle-input>

    <app-confirm-dialog-portal
      (onConfirm)="onAction($event)"
      [showModal]="showConfirmModal"
      [title]="!isActive ? 'Update doctor status' : 'Are you sure ?'"
      [description]="
        !isActive
          ? 'Doctor will be able to access his account and manage appointments.'
          : 'This will restrict doctor from accessing their account and managing appointments.'
      "
    >
    </app-confirm-dialog-portal> `,
})
export class DoctorStatusToggleComponent {
  @Input() isActive = false;

  showConfirmModal = false;

  constructor(private confirmPortalService: ConfirmDialogPortalService) {}

  onAction(reason: ConfirmPortalEventTypes) {
    if (reason == 'cancel') {
      this.showConfirmModal = false;
    } else {
      this.confirmPortalService.setLoadingState(true);

      setTimeout(() => {
        this.confirmPortalService.sendResponseStatus({
          message: 'Status updated successfully',
          type: 'success',
          timeOut: 5000,
          callback: () => {
            this.isActive = !this.isActive;
            this.showConfirmModal = false;
          },
        });
      }, 2000);
    }
  }
}
