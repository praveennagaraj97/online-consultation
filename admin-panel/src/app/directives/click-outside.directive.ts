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

  constructor(private _el: ElementRef<HTMLElement>) {}

  @HostListener('document:click', ['$event.target'])
  onclickOutside(ev: Element) {
    this.callback.emit(!this._el?.nativeElement?.contains(ev) || false);
  }
}
