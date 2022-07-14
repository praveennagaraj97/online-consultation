import { NgModule } from '@angular/core';
import { BlobImageViewerDirective } from './blob-image-viewer.directive';
import { ClickOutsideDirective } from './click-outside.directive';
import { IntersectionObserverDirective } from './intersection-observe.directive';
import { ScrollIntoViewDirective } from './scroll-into-view.directive';

@NgModule({
  declarations: [
    ClickOutsideDirective,
    IntersectionObserverDirective,
    ScrollIntoViewDirective,
    BlobImageViewerDirective,
  ],
  exports: [
    ClickOutsideDirective,
    IntersectionObserverDirective,
    ScrollIntoViewDirective,
    BlobImageViewerDirective,
  ],
})
export class DirectiveModule {}
