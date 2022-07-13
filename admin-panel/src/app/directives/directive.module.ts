import { NgModule } from '@angular/core';
import { ClickOutsideDirective } from './click-outside.directive';
import { IntersectionObserverDirective } from './intersection-observe.directive';

@NgModule({
  declarations: [ClickOutsideDirective, IntersectionObserverDirective],
  exports: [ClickOutsideDirective, IntersectionObserverDirective],
})
export class DirectiveModule {}
