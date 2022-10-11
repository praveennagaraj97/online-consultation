import { transition, trigger, useAnimation } from '@angular/animations';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { fadeInTransformAnimation } from 'src/app/animations';
import { ResponseMessageType } from 'src/app/types/app.types';

@Component({
  selector: 'app-message-tag',
  template: `
    <div
      [ngClass]="{
        'text-red-500': response?.type === 'error',
        'text-green-500': response?.type === 'success',
        'text-blue-500': response?.type === 'info'
      }"
      [class]="className"
    >
      <span *ngIf="response" @fadeIn>
        {{ response.message }}
      </span>
    </div>
  `,
  animations: [
    trigger('fadeIn', [
      transition('void => *', [useAnimation(fadeInTransformAnimation())]),
    ]),
  ],
})
export class MessageTagComponent {
  @Input() response?: ResponseMessageType | null = null;
  @Input() className: string = '';

  @Output() onEnd = new EventEmitter<void>();

  private timeId: any;
  ngOnChanges() {
    if (this.response) {
      console.log(this.response);
      clearTimeout(this.timeId);
      this.timeId = setTimeout(() => {
        this.onEnd.emit();
        this.response?.callback ? this.response?.callback() : {};
        this.response = null;
      }, this.response?.timeOut || 3000);
    }
  }
}
