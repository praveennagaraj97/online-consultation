import { animate, style, transition, trigger } from '@angular/animations';
import { Overlay, OverlayRef } from '@angular/cdk/overlay';
import { TemplatePortal } from '@angular/cdk/portal';
import {
  Component,
  Input,
  SimpleChanges,
  TemplateRef,
  ViewChild,
  ViewContainerRef,
} from '@angular/core';
import { Subscription } from 'rxjs';
import { clearSubscriptions } from 'src/app/utils/helpers';

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

  //   State
  private overlayRef?: OverlayRef;
  private templateRef?: TemplatePortal<HTMLDivElement>;

  constructor(
    private viewContainerRef: ViewContainerRef,
    private overlay: Overlay
  ) {}

  ngOnChanges(changes: SimpleChanges) {
    if (changes?.['showModal'].currentValue) {
      // Attach portal
      this.overlayRef?.attach(this.templateRef);
    } else if (!changes?.['showModal'].currentValue) {
      this.overlayRef?.dispose();
    }
  }

  ngAfterViewInit() {
    const overlay = this.overlay.create({ disposeOnNavigation: true });
    this.overlayRef = overlay;

    // Portal
    if (this.portalRef) {
      const portal = new TemplatePortal(this.portalRef, this.viewContainerRef);
      this.templateRef = portal;
    }

    // this.overlayRef?.attach(this.templateRef);

    this.overlayRef._outsidePointerEvents.subscribe({
      next: () => {
        this.overlayRef?.detach();
      },
    });
  }

  ngOnDestroy() {
    clearSubscriptions(this.subs$);
  }
}
