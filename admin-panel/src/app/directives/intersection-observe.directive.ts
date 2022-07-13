import { Directive, ElementRef, EventEmitter, Output } from '@angular/core';

@Directive({ selector: '[intersectionObserve]' })
export class IntersectionObserverDirective {
  @Output() callback: EventEmitter<boolean> = new EventEmitter();
  private observer?: IntersectionObserver;

  constructor(private _el: ElementRef) {}

  ngAfterViewInit() {
    const domElement = this._el.nativeElement;
    if (domElement) {
      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              this.callback.emit(true);
            } else {
              this.callback.emit(false);
            }
          });
        },
        { threshold: 0.7 }
      );

      observer.observe(domElement);
      this.observer = observer;
    }
  }

  ngOnDestroy() {
    this.observer?.disconnect();
  }
}
