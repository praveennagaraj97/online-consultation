import {
  Directive,
  ElementRef,
  EventEmitter,
  HostListener,
  Output,
} from '@angular/core';

@Directive({ selector: '[clickOutside]' })
export class ClickOutsideDirective {
  @Output() callback: EventEmitter<boolean> = new EventEmitter();

  constructor(private _el: ElementRef) {}

  @HostListener('document:click', ['$event.target'])
  onclickOutside(ev: EventTarget) {
    try {
      if (!this._el?.nativeElement.contains(ev as Element)) {
        this.callback.emit(false);
      } else {
        this.callback.emit(true);
      }
    } catch (error) {}
  }
}
