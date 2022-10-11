import { Injectable } from '@angular/core';
import { BehaviorSubject, Subject } from 'rxjs';
import { ResponseMessageType } from 'src/app/types/app.types';

@Injectable({ providedIn: 'any' })
export class ConfirmDialogPortalService {
  private isLoading$ = new BehaviorSubject<boolean>(false);
  private response$ = new Subject<ResponseMessageType>();

  //   Set Loading status for Ui to show
  setLoadingState(state: boolean) {
    this.isLoading$.next(state);
  }

  //   Listen to loading state observable
  get listenToLoadingState() {
    return this.isLoading$;
  }

  // When status is shared loading state will be togged after displaying the message
  sendResponseStatus(response: ResponseMessageType) {
    this.response$.next(response);
  }

  // Listen to response message
  get listenToResponse() {
    return this.response$;
  }
}
