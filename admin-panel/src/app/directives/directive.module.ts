import { NgModule } from '@angular/core';
import { BlobImageViewerDirective } from './blob-image-viewer.directive';
import { ClickOutsideDirective } from './click-outside.directive';
import { IntersectionObserverDirective } from './intersection-observe.directive';
import { PortalBackdropClickDirective } from './portal-backdrop-click.directive';
import { ScrollIntoViewDirective } from './scroll-into-view.directive';

@NgModule({
  declarations: [
    ClickOutsideDirective,
    IntersectionObserverDirective,
    ScrollIntoViewDirective,
    BlobImageViewerDirective,
    PortalBackdropClickDirective,
  ],
  exports: [
    ClickOutsideDirective,
    IntersectionObserverDirective,
    ScrollIntoViewDirective,
    BlobImageViewerDirective,
    PortalBackdropClickDirective,
  ],
})
export class DirectiveModule {}
