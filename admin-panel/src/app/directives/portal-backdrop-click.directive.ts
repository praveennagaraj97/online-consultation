import { OverlayRef } from '@angular/cdk/overlay';
import { Directive, EventEmitter, Input, Output } from '@angular/core';
import { Subscription } from 'rxjs';
import { clearSubscriptions } from '../utils/helpers';

@Directive({ selector: '[portalBackdropClick]' })
export class PortalBackdropClickDirective {
  // Subs
  private subs$: Subscription[] = [];

  @Input() overlayRef?: OverlayRef;
  @Output() callback: EventEmitter<boolean> = new EventEmitter();

  ngAfterViewInit() {
    if (this.overlayRef) {
      this.subs$.push(
        this.overlayRef.outsidePointerEvents().subscribe({
          next: () => {
            this.callback.emit(true);
          },
        })
      );
    } else {
      throw TypeError('Overlay ref is required');
    }
  }

  ngOnDestroy() {
    clearSubscriptions(this.subs$);
  }
}
