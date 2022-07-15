import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-email-sent',
  template: `
    <div class="container mx-auto">
      <div
        class="flex sm:dark:bg-gray-700 sm:bg-indigo-100 rounded-lg sm:shadow-xl py-6 items-center flex-wrap justify-center"
      >
        <img
          src="/assets/email-sent.png"
          alt="..."
          width="310"
          class="animate-pulse"
          height="530"
        />
        <div class="dark:text-gray-50 sm:p-6 p-2 sm:text-left text-center">
          <h1 class="text-lg font-semibold ">{{ title }}</h1>
          <p class="sm:text-sm text-xs mt-4 whitespace-pre">
            {{ description }}
          </p>
          <ng-content></ng-content>
        </div>
      </div>
    </div>
  `,
})
export class EmailSentComponent {
  @Input() title: string = '';
  @Input() description: string = '';
}
