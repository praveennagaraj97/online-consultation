import { transition, trigger, useAnimation } from '@angular/animations';
import { Component, ElementRef, Input, ViewChild } from '@angular/core';
import { AbstractControl, FormControl } from '@angular/forms';
import { fadeInTransformAnimation } from 'src/app/animations';
import { Country } from 'src/app/types/app.types';

@Component({
  selector: 'app-phone-input-component',
  templateUrl: 'phone-input.component.html',
  animations: [
    trigger('fadeIn', [
      transition('void => *', [useAnimation(fadeInTransformAnimation())]),
    ]),
  ],
})
export class PhoneInputComponent {
  // Props
  @Input() errors: { [key: string]: string } = {};
  @Input() fc!: AbstractControl<any, any> | undefined;
  @Input() fcPhoneCode!: AbstractControl<any, any> | undefined;
  @Input() showError = false;
  @Input() labelName = '';
  @Input() htmlFor = '';
  @Input() guideLine = '';
  @Input() placeholder = '';

  // Refs
  @ViewChild('countryCodeRef') countryCodeRef?: ElementRef<HTMLButtonElement>;

  // State
  selectedCountry: Country = {
    name: 'India',
    flag: '🇮🇳',
    code: 'IN',
    dial_code: '+91',
  };
  inputStyle: { [key: string]: string | number } = {};
  showCountryPicker = false;

  // Methods

  get control(): FormControl {
    return this.fc as FormControl;
  }

  get parseError(): string {
    const errorKey = Object.keys(this.fc?.errors || {})?.[0] || '';
    console.log(this.fc?.errors);
    return this.errors?.[errorKey] || 'Entered value is invalid';
  }

  private time_id: any;
  onCountryCodeSelect(country: Country) {
    clearTimeout(this.time_id);
    this.showCountryPicker = false;
    this.selectedCountry = country;
    this.fcPhoneCode?.setValue(country.dial_code);

    this.time_id = setTimeout(() => {
      const domRectRef =
        this.countryCodeRef?.nativeElement?.getBoundingClientRect();
      if (domRectRef) {
        this.inputStyle = {
          'padding-left.px': domRectRef.width + 8,
        };
      }
    }, 10);
  }
}
