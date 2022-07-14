import { NgModule } from '@angular/core';
import { ClickOutsideDirective } from './click-outside.directive';
import { IntersectionObserverDirective } from './intersection-observe.directive';
import { ScrollIntoViewDirective } from './scroll-into-view.directive';

@NgModule({
  declarations: [
    ClickOutsideDirective,
    IntersectionObserverDirective,
    ScrollIntoViewDirective,
  ],
  exports: [
    ClickOutsideDirective,
    IntersectionObserverDirective,
    ScrollIntoViewDirective,
  ],
})
export class DirectiveModule {}
