import { Directive, ElementRef, Input } from '@angular/core';

@Directive({ selector: '[scrollIntoView]' })
export class ScrollIntoViewDirective {
  constructor(private _el: ElementRef<HTMLElement>) {}
  @Input() viewOptions: ScrollIntoViewOptions = {
    behavior: 'smooth',
    block: 'center',
    inline: 'center',
  };

  ngAfterViewInit() {
    this._el.nativeElement.scrollIntoView(this.viewOptions);
  }
}
