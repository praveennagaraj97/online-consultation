import {
  animate,
  style,
  transition,
  trigger,
  useAnimation,
} from '@angular/animations';
import { Overlay, OverlayRef } from '@angular/cdk/overlay';
import { TemplatePortal } from '@angular/cdk/portal';
import {
  Component,
  EventEmitter,
  Input,
  Output,
  SimpleChanges,
  TemplateRef,
  ViewChild,
  ViewContainerRef,
} from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Subscription } from 'rxjs';
import { fadeInTransformAnimation } from 'src/app/animations';
import { ErrorResponse } from 'src/app/types/api.response.types';
import { APiResponseStatus } from 'src/app/types/app.types';
import { SpecialityEntity } from 'src/app/types/cms.response.types';
import { SpecialityFormDTO } from 'src/app/types/dto.types';
import { clearSubscriptions } from 'src/app/utils/helpers';
import {
  specialityErrors,
  SpecialityFormService,
} from './speciality-form.service';

@Component({
  selector: 'app-speciality-form-portal',
  templateUrl: 'speciality-form.component.html',
  animations: [
    trigger('swipeInOut', [
      transition('void => *', [
        style({ transform: 'translateX(100%)' }),
        animate('0.5s'),
      ]),
      transition('* => void', [
        style({ transform: 'translateX(0)' }),
        animate('0.5s', style({ transform: 'translateX(100%)' })),
      ]),
    ]),
    trigger('fadeIn', [
      transition('void => *', [useAnimation(fadeInTransformAnimation())]),
    ]),
  ],
})
export class SpecialityFormPortalComponent {
  // Refs
  @ViewChild('portalRef') portalRef?: TemplateRef<HTMLDivElement>;

  //   Subs
  private subs$: Subscription[] = [];

  //   Props
  @Input() showModal: boolean = false;
  @Input() formType: 'add' | 'edit' = 'add';

  // Event Emitters
  @Output() onBackdropClick = new EventEmitter<void>();
  @Output() onSuccessCallback = new EventEmitter();
  @Output() onAdded = new EventEmitter<SpecialityEntity>();

  //   State
  overlayRef?: OverlayRef;
  private templateRef?: TemplatePortal<HTMLDivElement>;
  errors = specialityErrors;
  showErrors = false;
  response: APiResponseStatus | null = null;
  submitting = false;

  // FormData
  thumbnail: File | null = null;
  specialityForm = new FormGroup<SpecialityFormDTO>({
    title: new FormControl('', { validators: [Validators.required] }),
  });

  constructor(
    private viewContainerRef: ViewContainerRef,
    private overlay: Overlay,
    private specialityService: SpecialityFormService
  ) {}

  ngOnChanges(changes: SimpleChanges) {
    if (changes?.['showModal'].currentValue) {
      this.overlayRef?.attach(this.templateRef);
    }
    if (!changes?.['showModal'].currentValue) {
      this.overlayRef?.detach();
    }
  }

  ngAfterViewInit() {
    const overlay = this.overlay.create({
      disposeOnNavigation: true,
    });
    this.overlayRef = overlay;

    // Portal
    if (this.portalRef) {
      const portal = new TemplatePortal(this.portalRef, this.viewContainerRef);
      this.templateRef = portal;
    }
  }

  handleOnSubmit() {
    if (this.specialityForm.invalid) {
      this.showErrors = true;
      return;
    }

    if (!this.thumbnail) {
      return this.setResponse({
        message: 'Thumbnail cannot be empty',
        type: 'error',
      });
    }

    this.submitting = true;

    this.subs$.push(
      this.specialityService
        .addNewSpeciality(this.specialityForm, this.thumbnail)
        .subscribe({
          next: (res) => {
            this.setResponse({ message: res.message, type: 'success' }, () => {
              this.onAdded.emit(res.result);
              this.submitting = false;
            });
          },
          error: (err: ErrorResponse) => {
            this.setResponse({ message: err.error.message, type: 'error' });
            this.submitting = false;
          },
        })
    );
  }

  handleThumbnailChange(file: File[]) {
    this.thumbnail = file[0];
  }

  private timeOutId: any;
  private setResponse(res: APiResponseStatus, callback?: () => void) {
    clearTimeout(this.timeOutId);
    this.response = res;

    setTimeout(() => {
      this.response = null;
      if (callback) {
        callback();
      }
    }, 3000);
  }

  ngOnDestroy() {
    clearSubscriptions(this.subs$);
  }
}
