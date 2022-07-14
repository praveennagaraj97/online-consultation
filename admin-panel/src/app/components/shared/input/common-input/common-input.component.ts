import { transition, trigger, useAnimation } from '@angular/animations';
import { Component, Input } from '@angular/core';
import { AbstractControl, FormControl } from '@angular/forms';
import { fadeInTransformAnimation } from 'src/app/animations';

@Component({
  selector: 'app-common-input-component',
  template: `
    <div>
      <label [for]="htmlFor" class="text-xs font-semibold">{{
        labelName
      }}</label>
      <small class="block text-xs mb-2 opacity-70" *ngIf="guideLine">
        {{ guideLine }}
      </small>
      <input
        [formControl]="control"
        [type]="type"
        name="htmlFor"
        autocomplete="off"
        class="common-input w-full input-focus input-colors p-2 rounded-md text-sm shadow-md"
        [placeholder]="placeholder"
        [ngClass]="
          showError
            ? fc?.invalid
              ? '!border-red-500/70'
              : '!border-green-500/70'
            : ''
        "
      />
      <div class="min-h-[18px] mt-0.5 smooth-animate">
        <small
          class="block text-red-500 text-xs"
          *ngIf="fc?.invalid && showError"
          @fadeIn
        >
          {{ parseError }}
        </small>
      </div>
    </div>
  `,
  animations: [
    trigger('fadeIn', [
      transition('void => *', [useAnimation(fadeInTransformAnimation())]),
    ]),
  ],
})
export class CommonInputComponent {
  @Input() errors: { [key: string]: string } = {};
  @Input() fc!: AbstractControl<any, any> | undefined;
  @Input() showError = false;
  @Input() labelName = '';
  @Input() htmlFor = '';
  @Input() guideLine = '';
  @Input() placeholder = '';
  @Input() type: string = 'text';

  get control(): FormControl {
    return this.fc as FormControl;
  }

  get parseError(): string {
    const errorKey = Object.keys(this.fc?.errors || {})?.[0] || '';
    console.log(this.fc?.errors, this.htmlFor);
    return this.errors?.[errorKey] || 'Entered value is invalid';
  }
}
