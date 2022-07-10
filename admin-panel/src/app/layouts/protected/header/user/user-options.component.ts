import { useAnimation } from '@angular/animations';
import { Component } from '@angular/core';
import { Subscription } from 'rxjs';
import { fadeInTransformAnimation } from 'src/app/animations';
import { ErrorResponse } from 'src/app/types/api.response.types';
import { UserEntity } from 'src/app/types/auth.response.types';
import { clearSubscriptions } from 'src/app/utils/helpers';
import { UserOptionsService } from './user-options.service';

@Component({
  selector: 'app-user-dropdown',
  templateUrl: 'user-options.component.html',
  animations: ['fadeIn', [useAnimation(fadeInTransformAnimation(400))]],
})
export class UserOptionsDropDownComponent {
  // Subscriptions
  private subs: Subscription[] = [];

  // State
  user: UserEntity | null = null;
  loading = false;

  constructor(private userService: UserOptionsService) {}

  ngOnInit() {
    this.loading = true;

    this.subs.push(
      this.userService.getProfileDetails().subscribe({
        next: (user) => {
          this.loading = false;
          this.user = user;
        },
        error: (err: ErrorResponse) => {
          this.loading = false;
          console.log(err);
        },
      })
    );
  }

  ngOnDestroy() {
    clearSubscriptions(this.subs);
  }
}
