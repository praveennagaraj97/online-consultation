import { Directive, ElementRef, Input } from '@angular/core';
import { convertFileUploadToBlobUrl } from '../utils/helpers';

@Directive({ selector: '[blobImageViewer]' })
export class BlobImageViewerDirective {
  @Input() file?: File;
  constructor(private _el: ElementRef<HTMLImageElement>) {}

  ngAfterViewInit() {
    if (this.file && this._el.nativeElement) {
      this._el.nativeElement.src = convertFileUploadToBlobUrl(this.file);
    }
  }
}
