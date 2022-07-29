import { Directive, ElementRef } from '@angular/core';

@Directive({ selector: '[stopPropagation]' })
export class StopPropagationDirective {
  constructor(private _el: ElementRef<HTMLElement>) {
    this._el.nativeElement.onclick = (ev) => {
      ev.stopPropagation();
      return ev;
    };
  }
}
