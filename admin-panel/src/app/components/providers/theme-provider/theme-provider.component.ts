import { animate, style, transition, trigger } from '@angular/animations';
import { Overlay, OverlayRef } from '@angular/cdk/overlay';
import {
  Component,
  TemplateRef,
  ViewChild,
  ViewContainerRef,
} from '@angular/core';

@Component({
  selector: 'app-theme-provider',
  templateUrl: 'theme-provider.component.html',
  animations: [
    trigger('backInRight', [
      transition('void => *', [
        style({ opacity: 0, transform: `translateX(100%)` }),
        animate('0.5s', style({ opacity: 1, transform: `translateX(0)` })),
      ]),
      transition('* => void', [
        style({ opacity: 1 }),
        animate('0.5s', style({ opacity: 0, transform: `translateX(100%)` })),
      ]),
    ]),
  ],
})
export class ThemeProviderComponent {
  @ViewChild('themePortal') themePortalRef?: TemplateRef<unknown>;
  selectedTheme: 'dark' | 'light' = 'dark';
  protected overlayRef?: OverlayRef;

  constructor(
    private overlay: Overlay,
    private viewContainerRef: ViewContainerRef
  ) {}

  ngOnInit() {
    // Custom Selected
    if (localStorage.getItem('theme')) {
      this.toggleTheme(localStorage.getItem('theme') as any);
      return;
    }

    // Old Browsers
    if (!window.matchMedia) {
      return;
    }

    const mql = window.matchMedia('(prefers-color-scheme: dark)');

    // System preference
    if (mql.matches) {
      this.toggleTheme('dark');
    } else {
      this.toggleTheme('light');
    }

    mql.addEventListener('change', () => {
      if (mql.matches) {
        this.toggleTheme('dark');
      } else {
        this.toggleTheme('light');
      }
    });
  }

  private toggleTheme(theme: 'dark' | 'light') {
    const dom = window.document.querySelector('body');
    dom?.removeAttribute('class');

    this.selectedTheme = theme;
    if (theme === 'dark') {
      dom?.classList.add('dark', 'theme-switch', 'dark-mode');
    } else {
      dom?.classList.add('theme-switch', 'light-mode');
    }
  }

  setTheme(theme: 'dark' | 'light') {
    localStorage.setItem('theme', theme);
    this.toggleTheme(theme);
  }

  closeSwitch() {
    this.overlayRef?.detach();
  }

  ngAfterViewInit() {
    const overlayRef = this.overlay.create({ disposeOnNavigation: false });
    if (this.themePortalRef) {
      // overlayRef.attach(
      //   new TemplatePortal(this.themePortalRef, this.viewContainerRef)
      // );
    }

    this.overlayRef = overlayRef;
  }
}
