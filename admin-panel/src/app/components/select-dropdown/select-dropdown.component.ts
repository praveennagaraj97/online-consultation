import { animate, style, transition, trigger } from '@angular/animations';
import { Overlay, OverlayRef } from '@angular/cdk/overlay';
import { TemplatePortal } from '@angular/cdk/portal';
import {
  Component,
  ElementRef,
  EventEmitter,
  Input,
  Output,
  TemplateRef,
  ViewChild,
  ViewContainerRef,
} from '@angular/core';
import { fromEvent, Subscription } from 'rxjs';
import { SelectOption } from 'src/app/types/app.types';
import { clearSubscriptions } from 'src/app/utils/helpers';

@Component({
  selector: 'app-select-dropdown',
  templateUrl: 'select-dropdown.component.html',
  animations: [
    trigger('openClose', [
      transition('void => *', [
        style({
          opacity: 0,
        }),
        animate('0.5s', style({ opacity: 1 })),
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
export class SelectDropdownComponent {
  // Subs
  private subs$: Subscription[] = [];

  // Ref
  @ViewChild('inputPosRef') inputPosRef?: ElementRef<HTMLDivElement>;
  @ViewChild('optionsRef') optionsRef?: TemplateRef<unknown>;

  private templateRef: TemplatePortal<unknown> | null = null;
  private overlayRef?: OverlayRef;

  // Props
  @Input() placeHolder = '';
  @Input() options: SelectOption[] = [];

  // Event Emitters
  @Output() onChange = new EventEmitter<SelectOption>();

  //   State
  selectedValue = '';
  showOption = false;
  renderPosition: { [key: string]: string | number } | null = null;

  constructor(
    private overlay: Overlay,
    private viewContainerRef: ViewContainerRef
  ) {}

  ngAfterViewInit() {
    const domRef = this.optionsRef;
    if (!domRef) {
      return;
    }

    const overlay = this.overlay.create({ disposeOnNavigation: true });
    this.templateRef = new TemplatePortal(domRef, this.viewContainerRef!);
    this.overlayRef = overlay;
  }

  private listenToWindowChanges() {
    this.subs$.push(
      fromEvent(window, 'resize').subscribe({
        next: () => {
          if (this.showOption) {
            this.updateModalPosition();
          }
        },
      })
    );

    this.subs$.push(
      fromEvent(window, 'scroll').subscribe({
        next: () => {
          if (this.showOption) {
            this.updateModalPosition();
          }
        },
      })
    );
  }

  private updateModalPosition() {
    const positionDom = this.inputPosRef?.nativeElement;
    if (!positionDom) {
      return;
    }

    const clientRect = positionDom.getBoundingClientRect();
    this.renderPosition = {
      'left.px': clientRect.left,
      'top.px': clientRect.top + clientRect.height,
      'width.px': clientRect.width,
    };
  }

  showOptions() {
    this.overlayRef?.attach(this.templateRef);
    this.updateModalPosition();
    this.listenToWindowChanges();
    this.showOption = true;
  }

  selectAnOption(option: SelectOption) {
    this.selectedValue = option.title;
    this.onChange.emit(option);
    this.hideOptions(true);
  }

  hideOptions(state: boolean) {
    if (state) {
      this.overlayRef?.detach();
      this.showOption = false;
      this.renderPosition = null;
      clearSubscriptions(this.subs$);
    }
  }

  ngOnDestroy() {
    clearSubscriptions(this.subs$);
  }
}
