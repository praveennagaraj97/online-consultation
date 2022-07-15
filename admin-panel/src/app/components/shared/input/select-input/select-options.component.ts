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
import { debounceTime, fromEvent, Subject, Subscription } from 'rxjs';
import { SelectOption } from 'src/app/types/app.types';
import { clearSubscriptions } from 'src/app/utils/helpers';

@Component({
  selector: 'app-select-input-options',
  template: `
    <ng-template #portalRef>
      <div
        (click)="$event.stopPropagation()"
        class="dark:bg-gray-700 bg-gray-50 rounded-lg fixed z-10 shadow-2xl overflow-hidden"
        [ngStyle]="posStyle"
        @openClose
      >
        <div class="relative" *ngIf="isAsync">
          <div class="absolute left-2 top-0 bottom-0 flex items-center">
            <app-search-icon
              className="w-4 h-4 dark:fill-gray-300 dark:stroke-gray-300 stroke-gray-700 fill-gray-700 pt-1"
            ></app-search-icon>
          </div>
          <input
            type="text"
            class="w-full p-2 pl-8 common-input text-sm input-focus input-colors !border-t-0 !border-x-0 rounded-lg mt-1"
            [placeholder]="placeholder"
            [ngModel]="searchFormValue"
            (ngModelChange)="onSearchChange($event)"
          />
        </div>

        <!-- Multiple Inputs Display -->

        <div class="max-h-[144px] overflow-y-auto inner-scrollbar">
          <button
            class="p-2  smooth-animate hover:bg-sky-500/30 w-full  flex space-x-2 items-center text-xs"
            *ngFor="let option of options"
            [title]="option.title"
            (click)="handleOnSelect(option)"
          >
            {{ option.title }}
          </button>
          <div
            *ngIf="hasMore && !loading && isAsync"
            intersectionObserve
            (callback)="loadMoreEvent($event)"
          ></div>
          <app-spinner-icon
            *ngIf="loading"
            className="animate-spin dark:fill-gray-300 dark:stroke-gray-300 stroke-gray-700 fill-gray-700 w-6 h-6 mx-auto my-3 py-1"
          ></app-spinner-icon>
        </div>
      </div>
    </ng-template>
  `,
  animations: [
    trigger('openClose', [
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
export class SelectInputOptionsComponent {
  // Subs
  private subs$: Subscription[] = [];

  // Props
  @Input() renderPosition: DOMRect | null = null;
  @Input() showOptions = false;
  @Input() options: SelectOption[] = [];
  @Input() placeholder = '';
  @Input() loading = false;
  @Input() hasMore = false;
  @Input() isAsync = false;

  // Event Emitters
  @Output() onChange = new EventEmitter<SelectOption>(false);
  @Output() loadMore = new EventEmitter<void>(false);
  @Output() onSearch = new EventEmitter<string>(false);
  debouncer: Subject<string> = new Subject<string>();

  // State
  dummyCards = new Array(10).fill('');
  posStyle: { [key: string]: string | number } = {};
  searchFormValue = '';

  //   Refs
  private overlayRef?: OverlayRef;
  @ViewChild('portalRef') portalRef?: TemplateRef<unknown>;
  private templateRef!: TemplatePortal;

  constructor(
    private overlay: Overlay,
    private viewContainer: ViewContainerRef
  ) {}

  ngOnInit() {
    this.subs$.push(
      this.debouncer.pipe(debounceTime(250)).subscribe({
        next: (value) => {
          this.onSearch?.emit(value);
        },
      })
    );
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes?.['showOptions']?.currentValue) {
      this.overlayRef?.attach(this.templateRef);
      this.updateModalPosition();
      this.listenToWindowChanges();
    }

    if (changes?.['showOptions']?.currentValue === false) {
      this.onSearchChange('');
      this.overlayRef?.detach();
    }
  }

  ngAfterViewInit() {
    if (!this.portalRef) {
      return;
    }

    const overlay = this.overlay.create({ disposeOnNavigation: true });
    this.overlayRef = overlay;

    this.templateRef = new TemplatePortal(this.portalRef, this.viewContainer);
  }

  handleOnSelect(opt: SelectOption) {
    this.onChange.emit(opt);
  }

  private updateModalPosition() {
    const domPos = this.renderPosition;
    if (!domPos) {
      return;
    }

    this.posStyle = {
      'top.px': domPos.top + domPos.height,
      'left.px': domPos.left,
      'width.px': domPos.width,
    };
  }

  private listenToWindowChanges() {
    this.subs$.push(
      fromEvent(window, 'resize').subscribe({
        next: () => {
          this.updateModalPosition();
        },
      })
    );

    this.subs$.push(
      fromEvent(window, 'scroll').subscribe({
        next: () => {
          this.updateModalPosition();
        },
      })
    );
  }

  loadMoreEvent(state: boolean) {
    if (state) {
      this.loadMore?.emit();
    }
  }

  onSearchChange(term: string) {
    this.debouncer.next(term);
  }

  ngOnDestroy() {
    clearSubscriptions(this.subs$);
  }
}
