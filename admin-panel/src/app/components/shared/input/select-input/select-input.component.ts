import { transition, trigger, useAnimation } from '@angular/animations';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { AbstractControl, FormControl } from '@angular/forms';
import { fadeInTransformAnimation } from 'src/app/animations';
import { SelectOption } from 'src/app/types/app.types';

@Component({
  selector: 'app-select-input-component',
  templateUrl: 'select-input.component.html',
  animations: [
    trigger('fadeIn', [
      transition('void => *', [useAnimation(fadeInTransformAnimation())]),
    ]),
  ],
})
export class SelectInputComponent {
  // Props
  @Input() errors: { [key: string]: string } = {};
  @Input() fc!: AbstractControl<any, any> | undefined;
  @Input() showError = false;
  @Input() labelName = '';
  @Input() htmlFor = '';
  @Input() guideLine = '';
  @Input() placeholder = '';
  @Input() type: string = 'text';
  @Input() isAsync = false;
  @Input() hasNext = false;
  @Input() isLoading = false;
  @Input() options: SelectOption[] = [];
  @Input() searchPlaceholder = '';
  @Input() isMulti = false;
  @Input() multiInputIgnoreKeys: string[] = [];
  multipleOptions: SelectOption[] = [];

  // Event Emitters
  @Output() loadMore = new EventEmitter<void>(false);
  @Output() onSearch = new EventEmitter<string>(false);
  @Output() onChange = new EventEmitter<SelectOption>(false);

  // State
  showOptions = false;
  selectOption = '';
  multiViewScrollViewOptions: ScrollIntoViewOptions = {
    behavior: 'smooth',
    block: 'nearest',
    inline: 'start',
  };

  get control(): FormControl {
    return this.fc as FormControl;
  }

  get parseError(): string {
    const errorKey = Object.keys(this.fc?.errors || {})?.[0] || '';
    console.log(this.fc?.errors);
    return this.errors?.[errorKey] || 'Entered value is invalid';
  }

  onSelect(opt: SelectOption) {
    if (this.isMulti) {
      // Custom event trigger keys
      if (this.multiInputIgnoreKeys.includes(opt.value)) {
        this.showOptions = false;
        return;
      }
      if (this.multipleOptions.find((mo) => mo.value == opt.value)) {
        this.multipleOptions = this.multipleOptions.filter(
          (option) => option.value != opt.value
        );
      } else {
        this.multipleOptions = [...this.multipleOptions, opt];
      }

      this.onChange.emit(opt);
      return;
    }

    this.showOptions = false;
    this.selectOption = opt.title;
    this.onChange.emit(opt);
  }

  // Remove multi select option
  removeSelection(opt: SelectOption) {
    this.multipleOptions = this.multipleOptions.filter(
      (option) => option.value != opt.value
    );
    this.onChange.emit(opt);
  }
}
