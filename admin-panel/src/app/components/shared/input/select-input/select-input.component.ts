import { transition, trigger, useAnimation } from '@angular/animations';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { AbstractControl, FormControl } from '@angular/forms';
import { fadeInTransformAnimation } from 'src/app/animations';
import { SelectOption } from 'src/app/types/app.types';

@Component({
  selector: 'app-select-input',
  templateUrl: 'select-input.component.html',
  animations: [
    trigger('fadeIn', [
      transition('void => *', [useAnimation(fadeInTransformAnimation())]),
    ]),
  ],
})
export class SelectInputComponent {
  // Props
  // Errors Messages for this input
  @Input() errors: { [key: string]: string } = {};
  // FormControl on which input will error existence will be checked
  @Input() fc?: AbstractControl<any, any> | undefined;
  // COndition for when to display the errors
  @Input() showError = false;
  // Label name for input - Optional
  @Input() labelName = '';
  // Optional
  @Input() htmlFor = '';
  // Optional
  @Input() guideLine = '';
  // Optional
  @Input() placeholder = '';
  // Optional
  @Input() type: string = 'text';
  // Optional
  @Input() isAsync = false;
  // Optional
  @Input() hasNext = false;
  // Optional
  @Input() isLoading = false;
  @Input() options: SelectOption[] = [];
  // Optional
  @Input() searchPlaceholder = '';
  // Optional
  @Input() isMulti = false;
  // Optional
  @Input() multiInputIgnoreKeys: string[] = [];
  // Optional
  @Input() defaultValue: string | null = null;

  // State
  multipleOptions: SelectOption[] = [];

  // Event Emitters
  @Output() loadMore = new EventEmitter<void>(false);
  @Output() onSearch = new EventEmitter<string>(false);
  @Output() onChange = new EventEmitter<SelectOption>(false);

  // State
  showOptions = false;
  selectOptionValue = '';
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

    return this.errors?.[errorKey] || 'Entered value is invalid';
  }

  ngOnInit() {
    if (this.defaultValue) {
      this.selectOptionValue = this.defaultValue;
    }
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

      this.onChange?.emit(opt);
      return;
    }

    this.showOptions = false;
    this.selectOptionValue = opt.title;
    this.onChange?.emit(opt);
  }

  // Remove multi select option
  removeSelection(opt: SelectOption) {
    this.multipleOptions = this.multipleOptions.filter(
      (option) => option.value != opt.value
    );
    this.onChange?.emit(opt);
  }
}
