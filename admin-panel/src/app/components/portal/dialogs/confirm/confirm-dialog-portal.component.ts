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
import {
  ConfirmPortalEventTypes,
  ResponseMessageType,
} from 'src/app/types/app.types';
import { clearSubscriptions } from 'src/app/utils/helpers';
import { ConfirmDialogPortalService } from './confirm-dialog-portal.service';

@Component({
  selector: 'app-confirm-dialog-portal',
  templateUrl: 'confirm-dialog-portal.component.html',
  animations: [
    trigger('swipeInOut', [
      transition('void => *', [
        style({ opacity: 0 }),
        animate('0.5s', style({ opacity: 1 })),
      ]),
    ]),
  ],
})
export class ConfirmDialogPortalComponent {
  // Subs
  private subs$: Subscription[] = [];

  // State
  isLoading = false;
  response: ResponseMessageType | null = null;

  // Refs
  @ViewChild('portalRef') portalRef?: TemplateRef<HTMLDivElement>;
  private overlayRef!: OverlayRef;
  private templatePortal?: TemplatePortal<HTMLDivElement>;

  // Events
  @Output() onConfirm = new EventEmitter<ConfirmPortalEventTypes>();
  @Input() showModal = false;

  // Props
  @Input() title?: string = 'Are you sure you want to continue ?';
  @Input() description?: string = '';
  @Input() cancelName?: string = 'Cancel';
  @Input() confirmName?: string = 'Confirm';

  constructor(
    private overlay: Overlay,
    private viewContainerRef: ViewContainerRef,
    private confirmPortalService: ConfirmDialogPortalService
  ) {}

  ngAfterViewInit() {
    if (!this.portalRef) {
      return;
    }

    this.overlayRef = this.overlay.create({
      disposeOnNavigation: true,
    });

    this.templatePortal = new TemplatePortal<HTMLDivElement>(
      this.portalRef,
      this.viewContainerRef
    );
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes?.['showModal']?.currentValue) {
      this.overlayRef?.attach(this.templatePortal);
    } else if (changes?.['showModal']?.currentValue == false) {
      this.cleanUp();
    }
  }

  onClickOutSide(state: boolean) {
    if (this.isLoading) {
      return;
    }

    if (state) {
      this.onConfirm.emit('cancel');
    }
  }

  onConfirmClick() {
    this.onConfirm.emit('confirm');
    this.subs$.push(
      this.confirmPortalService.listenToLoadingState.subscribe({
        next: (loading) => {
          this.isLoading = loading;
        },
      })
    );

    this.subs$.push(
      this.confirmPortalService.listenToLoadingState.subscribe({
        next: (loading) => {
          this.isLoading = loading;
        },
      })
    );

    this.subs$.push(
      this.confirmPortalService.listenToResponse.subscribe({
        next: (res) => {
          this.response = res;
        },
      })
    );
  }

  private cleanUp() {
    clearSubscriptions(this.subs$);
    this.response = null;
    this.isLoading = false;
    this.overlayRef?.detach();
  }

  ngOnDestroy() {
    this.overlayRef?.dispose();
  }
}
