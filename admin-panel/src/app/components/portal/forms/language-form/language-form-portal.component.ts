import { animate, style, transition, trigger } from '@angular/animations';
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
import { languageErrors } from 'src/app/errors/cms-form.errors';
import { ErrorResponse } from 'src/app/types/api.response.types';
import { APiResponseStatus } from 'src/app/types/app.types';
import { LanguageEntity } from 'src/app/types/cms.response.types';
import { clearSubscriptions } from 'src/app/utils/helpers';
import { LanguageFormService } from './language-form.service';

@Component({
  selector: 'app-language-form-portal',
  templateUrl: 'language-form.component.html',
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
  ],
})
export class LanguageFormPortalComponent {
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
  @Output() onAdded = new EventEmitter<LanguageEntity>();

  //   State
  overlayRef?: OverlayRef;
  private templateRef?: TemplatePortal<HTMLDivElement>;
  errors = languageErrors;
  showErrors = false;
  response: APiResponseStatus | null = null;
  submitting = false;

  languageForm = new FormGroup({
    name: new FormControl('', { validators: [Validators.required] }),
    locale_name: new FormControl('', { validators: [Validators.required] }),
  });

  constructor(
    private viewContainerRef: ViewContainerRef,
    private overlay: Overlay,
    private langFormService: LanguageFormService
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
    if (this.languageForm.invalid) {
      this.showErrors = true;
      return;
    }

    this.submitting = true;

    this.subs$.push(
      this.langFormService.addNewLanguage(this.languageForm).subscribe({
        next: (res) => {
          this.setResponse({ message: res.message, type: 'success' }, () => {
            this.submitting = false;
            this.onAdded.emit(res.result);
            this.languageForm.reset();
          });
        },
        error: ({ error }: ErrorResponse) => {
          this.setResponse({ message: error.message, type: 'error' });
          this.submitting = false;
        },
      })
    );
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
