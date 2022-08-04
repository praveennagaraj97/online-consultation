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
import { Subscription } from 'rxjs';
import { ConfirmPortalEventTypes } from 'src/app/types/app.types';

@Component({
  selector: 'app-confirm-dialog-portal',
  templateUrl: 'confirm-dialog-portal.component.html',
  animations: [
    trigger('swipeInOut', [
      transition('void => *', [
        style({ transform: 'translateY(-200%)', opacity: 0 }),
        animate('0.5s', style({ transform: 'translateY(80px)', opacity: 1 })),
      ]),
      transition('* => void', [
        style({ transform: 'translateY(80px)', opacity: 1 }),
        animate('0.5s', style({ transform: 'translateY(-200%)', opacity: 0 })),
      ]),
    ]),
  ],
})
export class ConfirmDialogPortalComponent {
  // Subs
  private subs$: Subscription[] = [];

  // Refs
  @ViewChild('portalRef') portalRef?: TemplateRef<HTMLDivElement>;
  overlayRef!: OverlayRef;
  private templatePortal?: TemplatePortal<HTMLDivElement>;

  // Events
  @Output() onConfirm = new EventEmitter<ConfirmPortalEventTypes>();
  @Input() showModal = false;

  constructor(
    private overlay: Overlay,
    private viewContainerRef: ViewContainerRef
  ) {}

  ngAfterViewInit() {
    if (!this.portalRef) {
      return;
    }

    this.overlayRef = this.overlay.create({
      disposeOnNavigation: true,
    });

    this.subs$.push(
      this.overlayRef.outsidePointerEvents().subscribe({
        next: () => {
          this.onConfirm.emit('cancel');
        },
      })
    );

    this.templatePortal = new TemplatePortal<HTMLDivElement>(
      this.portalRef,
      this.viewContainerRef
    );
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes?.['showModal'].currentValue) {
      this.overlayRef?.attach(this.templatePortal);
    }

    if (!changes?.['showModal'].currentValue) {
      this.overlayRef?.detach();
    }
  }

  onClickOutSide(state: boolean) {
    if (!state) {
      this.onConfirm.emit('cancel');
    }
  }

  ngOnDestroy() {
    this.overlayRef?.dispose();
  }
}
