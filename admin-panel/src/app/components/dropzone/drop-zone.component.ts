import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-dropzone',
  template: `
    <div
      class="w-full h-full"
      (drop)="handleOnDrop($event)"
      (dragover)="$event.preventDefault(); $event.stopPropagation()"
      (click)="inputRef.click()"
    >
      <ng-content></ng-content>
      <input
        type="file"
        [accept]="accept"
        (change)="handleOnInput($event)"
        (click)="resetFileInput($event)"
        hidden
        #inputRef
      />
    </div>
  `,
})
export class DropzoneComponent {
  @Input() isMulti = false;
  @Input() dropfileType = 'image';
  @Input() accept = 'image/*';

  // List for single file
  @Output() onChange = new EventEmitter<File[]>();

  handleOnDrop(ev: DragEvent) {
    ev.preventDefault();

    ev.preventDefault();
    const dataTransfer = ev.dataTransfer;
    if (dataTransfer) {
      this.emitFiles(dataTransfer.files);
    }
  }

  resetFileInput(ev: Event) {
    const element = ev.currentTarget as HTMLInputElement;
    element.value = '';
  }

  handleOnInput(ev: Event) {
    ev.preventDefault();

    const element = ev.currentTarget as HTMLInputElement;
    let fileList: FileList | null = element.files;
    if (fileList) {
      this.emitFiles(fileList);
    }
  }

  private emitFiles(filesList: FileList) {
    if (this.isMulti) {
      const files: File[] = [];
      for (let i = 0; i < filesList.length; i++) {
        const file = filesList.item(i);
        if (file && file.type.includes(this.dropfileType)) {
          files.push(file);
        }
      }

      this.onChange.emit(files);
    } else {
      const file = filesList.item(0);

      if (file && file.type.includes(this.dropfileType)) {
        this.onChange.emit([file]);
      }
    }
  }
}
