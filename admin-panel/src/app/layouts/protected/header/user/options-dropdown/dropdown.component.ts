import { animate, style, transition, trigger } from '@angular/animations';
import { Overlay, OverlayRef } from '@angular/cdk/overlay';
import { TemplatePortal } from '@angular/cdk/portal';
import {
  Component,
  EventEmitter,
  Input,
  Output,
  TemplateRef,
  ViewChild,
  ViewContainerRef,
} from '@angular/core';

@Component({
  selector: 'app-user-options-dropdown',
  templateUrl: 'dropdown.component.html',
  animations: [
    trigger('fadeInOut', [
      transition('void => *', [
        style({
          opacity: 0,
          transform: 'translateY(10px)',
        }),
        animate('0.5s', style({ opacity: 1, transform: 'translateY(0px)' })),
      ]),
      transition('* => void', [
        style({
          opacity: 1,
        }),
        animate('0.5s', style({ opacity: 0 })),
      ]),
    ]),
  ],
})
export class UserOptionsDropdownComponent {
  // Props
  @Input() position: DOMRect | null = null;

  // Refs
  @ViewChild('dropdownRef') dropdownRef?: TemplateRef<HTMLDivElement>;
  private overlayRef?: OverlayRef;

  // State
  domPosition: { top: string; left: string } = { left: '', top: '' };

  // Emitters
  @Output() onClose = new EventEmitter<boolean>(false);

  constructor(
    private overlay: Overlay,
    private viewContainerRef: ViewContainerRef
  ) {}

  ngAfterViewInit() {
    if (this.dropdownRef) {
      const overlay = this.overlay.create({
        disposeOnNavigation: true,
      });

      this.domPosition = {
        left: `${(this.position?.left || 0) + 15}px`,
        top: `${(this.position?.top || 0) + 55}px`,
      };

      overlay.attach(
        new TemplatePortal(this.dropdownRef, this.viewContainerRef)
      );

      this.overlayRef = overlay;
    }
  }

  onCloseCallback(state: boolean) {
    if (state) {
      this.overlayRef?.detach();
      this.onClose.emit(true);
    }
  }
}
