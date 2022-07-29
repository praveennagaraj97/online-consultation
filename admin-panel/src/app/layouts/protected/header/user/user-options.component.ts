import { transition, trigger, useAnimation } from '@angular/animations';
import { Component, ElementRef, ViewChild } from '@angular/core';
import { Subscription } from 'rxjs';
import { fadeInTransformAnimation } from 'src/app/animations';
import { ErrorResponse } from 'src/app/types/api.response.types';
import { UserEntity } from 'src/app/types/auth.response.types';
import { clearSubscriptions } from 'src/app/utils/helpers';
import { UserOptionsService } from './user-options.service';

@Component({
  selector: 'app-user-options',
  template: `
    <div
      class="flex items-center"
      role="button"
      *ngIf="!loading && user"
      @fadeIn
      #optionsRef
      (click)="handleshowDropdown($event)"
    >
      <div class="mr-2">
        <p class="text-sm leading-normal font-semibold">{{ user?.name }}</p>
        <small class="text-right text-xs block capitalize">
          {{ user.role | replaceUnderscore }}
        </small>
      </div>

      <app-user-rounded-icon
        className="w-10 h-10 dark:fill-gray-50 dark:stroke-gray-50"
      ></app-user-rounded-icon>
    </div>

    <app-profile-skeleton *ngIf="loading"></app-profile-skeleton>

    <!-- Dropdown Portal -->
    <app-user-options-dropdown
      [showDropdown]="showDropdown"
      [position]="domRectOptions"
      (onClose)="showDropdown = false"
    ></app-user-options-dropdown>
  `,
  animations: [
    trigger('fadeIn', [
      transition('void => *', [useAnimation(fadeInTransformAnimation())]),
    ]),
  ],
})
export class UserOptionsComponent {
  // Subscriptions
  private subs$: Subscription[] = [];

  // State
  user: UserEntity | null = null;
  loading = false;
  showDropdown = false;
  domRectOptions: DOMRect | null = null;

  // Ref
  @ViewChild('optionsRef') optionsRef?: ElementRef<HTMLDivElement>;

  constructor(private userService: UserOptionsService) {}

  ngOnInit() {
    this.loading = true;

    this.subs$.push(
      this.userService.getProfileDetails().subscribe({
        next: (user) => {
          this.loading = false;
          this.user = user;
        },
        error: (err: ErrorResponse) => {
          this.loading = false;
          alert(err);
        },
      })
    );
  }

  handleshowDropdown(ev: MouseEvent) {
    ev.stopPropagation();
    const domRef = this.optionsRef?.nativeElement;
    if (domRef) {
      this.domRectOptions = domRef.getBoundingClientRect();
      this.showDropdown = true;
    }
  }

  ngOnDestroy() {
    clearSubscriptions(this.subs$);
  }
}
