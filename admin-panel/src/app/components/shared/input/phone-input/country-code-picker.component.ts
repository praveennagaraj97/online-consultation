import { animate, style, transition, trigger } from '@angular/animations';
import { Overlay, OverlayRef } from '@angular/cdk/overlay';
import { TemplatePortal } from '@angular/cdk/portal';
import { HttpClient } from '@angular/common/http';
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
import { fromEvent, Subscription } from 'rxjs';
import { additionalRoutes } from 'src/app/api-routes/routes';
import { Country } from 'src/app/types/app.types';
import { clearSubscriptions } from 'src/app/utils/helpers';

@Component({
  selector: 'app-country-code-picker',
  template: `
    <ng-template #portalRef>
      <div
        (click)="$event.stopPropagation()"
        class="dark:bg-gray-700 bg-gray-50 rounded-lg fixed z-10 shadow-2xl overflow-hidden"
        [ngStyle]="posStyle"
        @openClose
        portalBackdropClick
        [overlayRef]="overlayRef"
        (callback)="onBackdropClose.emit()"
      >
        <input
          [ngModel]="searchFormValue"
          (ngModelChange)="onSearchChange($event)"
          type="text"
          class="w-full py-2 px-1 common-input text-sm input-focus input-colors !border-t-0 !border-x-0 rounded-lg mt-1"
          placeholder="Search your country"
        />
        <div class="max-h-[260px] overflow-y-auto inner-scrollbar">
          <button
            (click)="onSelect(country)"
            class="p-2  smooth-animate hover:bg-sky-500/30 w-full  flex space-x-2 items-center"
            *ngFor="let country of countries"
            [title]="country.name"
          >
            <p class="pt-0.5 text-lg">
              {{ country.flag }}
            </p>
            <p class="text-sm cut-text-1">
              {{ country.name }}
            </p>
          </button>
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
export class CountryCodePickerComponent {
  // Props
  @Input() renderPosition: DOMRect | null = null;
  @Input() showCountryPicker = false;
  @Output() onChange = new EventEmitter<Country>(false);
  @Output() onBackdropClose = new EventEmitter<void>();

  // State
  countries: Country[] = [];
  private countriesData: Country[] = [];
  posStyle: { [key: string]: string | number } = {};
  searchFormValue = '';

  //   Refs
  overlayRef?: OverlayRef;
  @ViewChild('portalRef') portalRef?: TemplateRef<unknown>;
  private templateRef!: TemplatePortal;

  // Subsriptions
  private subs$: Subscription[] = [];

  constructor(
    private http: HttpClient,
    private overlay: Overlay,
    private viewContainer: ViewContainerRef
  ) {}

  ngOnInit() {
    this.getCountriesCodes();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes?.['showCountryPicker']?.currentValue) {
      this.overlayRef?.attach(this.templateRef);

      this.updateModalPosition();
      this.listenToWindowChanges();
    }

    if (changes?.['showCountryPicker']?.currentValue === false) {
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

  private getCountriesCodes() {
    this.subs$.push(
      this.http.get<Country[]>(additionalRoutes.GetCountries).subscribe({
        next: (countries) => {
          this.countries = countries;
          this.countriesData = countries;
        },
        error: (err) => {
          alert('Something went wrong while loading countries');
        },
      })
    );
  }

  onSelect(country: Country) {
    this.countries = this.countriesData;
    this.onChange.emit(country);
  }

  onSearchChange(value: string) {
    const searchedValue = value.toLowerCase();
    this.countries = this.countriesData;

    const searchedCountries = this.countries.filter((country) => {
      const countryNameLower = country.name.toLowerCase();
      return (
        countryNameLower === searchedValue ||
        countryNameLower.startsWith(searchedValue) ||
        countryNameLower.includes(searchedValue)
      );
    });

    this.countries = searchedCountries;
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

  ngOnDestroy() {
    clearSubscriptions(this.subs$);
  }
}
