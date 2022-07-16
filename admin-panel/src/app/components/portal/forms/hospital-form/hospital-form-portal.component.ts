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
import { hospitalFormErrors } from 'src/app/errors/hospital-form.errors';
import { SelectOption } from 'src/app/types/app.types';
import { HospitalFormDTO } from 'src/app/types/dto.types';
import { clearSubscriptions } from 'src/app/utils/helpers';
import { HospitalFormService } from './hospital-form.service';

@Component({
  selector: 'app-hospital-form-portal',
  templateUrl: 'hospital-form.component.html',
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
export class HospitalFormPortalComponent {
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

  //   State
  overlayRef?: OverlayRef;
  private templateRef?: TemplatePortal<HTMLDivElement>;
  errors = hospitalFormErrors;
  countriesOptions: SelectOption[] = [];
  showErrors = false;

  constructor(
    private viewContainerRef: ViewContainerRef,
    private overlay: Overlay,
    private hspFormService: HospitalFormService
  ) {}

  hospitalForm = new FormGroup<HospitalFormDTO>({
    name: new FormControl('', { validators: Validators.required }),
    country: new FormControl('India', { validators: Validators.required }),
    city: new FormControl('', { validators: Validators.required }),
    address: new FormControl('', { validators: Validators.required }),
  });

  ngOnChanges(changes: SimpleChanges) {
    if (changes?.['showModal'].currentValue) {
      this.overlayRef?.attach(this.templateRef);
    }
    if (!changes?.['showModal'].currentValue) {
      this.overlayRef?.detach();
    }
  }

  ngOnInit() {
    this.getCountries();
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

  private getCountries() {
    this.subs$.push(
      this.hspFormService.getCountries().subscribe({
        next: (res) => {
          this.countriesOptions = res;
        },
        error: (err) => {
          alert(err);
        },
      })
    );
  }

  handleSubmit() {
    if (this.hospitalForm.invalid) {
      this.showErrors = true;
      return;
    }

    this.subs$.push(
      this.hspFormService.addNewHospital(this.hospitalForm).subscribe({
        next: (res) => {
          console.log(res);
        },
        error: (err) => {
          console.log(err);
        },
      })
    );
  }

  ngOnDestroy() {
    clearSubscriptions(this.subs$);
  }
}
