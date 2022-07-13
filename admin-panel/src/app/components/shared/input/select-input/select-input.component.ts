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

  @Output() loadMore = new EventEmitter<void>(false);
  @Output() onSearch = new EventEmitter<string>(false);
  @Output() onChange = new EventEmitter<SelectOption>(false);

  // State
  showOptions = false;
  selectOption = '';

  get control(): FormControl {
    return this.fc as FormControl;
  }

  get parseError(): string {
    const errorKey = Object.keys(this.fc?.errors || {})?.[0] || '';
    console.log(this.fc?.errors);
    return this.errors?.[errorKey] || 'Entered value is invalid';
  }

  onSelect(opt: SelectOption) {
    this.selectOption = opt.title;
    this.onChange.emit(opt);
  }
}
